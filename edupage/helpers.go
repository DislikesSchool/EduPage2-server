package edupage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

var (
	ErrorUnitialized = errors.New("unitialized")
	ErrorNotFound    = errors.New("not found")
)

// GetStudentID is used to retrieve the client's student ID.
// Returns ErrorUnitialized if the user object hasn't been initialized.
func (client *EdupageClient) GetStudentID() (string, error) {
	if client.user == nil {
		return "", ErrorUnitialized
	}
	return client.user.UserRow.StudentID, nil
}

// GetSubjectByID is used to retrieve the subject by it's specified ID.
// Returns ErrorNotFound if the subject can't be found.
// Returns ErrorUnitialized if the user object hasn't been initialized.
func (client *EdupageClient) GetSubjectByID(id string) (model.Subject, error) {
	if client.user == nil {
		return model.Subject{}, ErrorUnitialized
	}

	if teacher, ok := client.user.DBI.Subjects[id]; ok {
		return teacher, nil
	}
	return model.Subject{}, ErrorNotFound
}

// GetTeacherByID is used to retrieve the teacher by their specified ID.
// Returns ErrorNotFound if the teacher can't be found.
// Returns ErrorUnitialized if the user object hasn't been initialized.
func (client *EdupageClient) GetTeacherByID(id string) (model.Teacher, error) {
	if client.user == nil {
		return model.Teacher{}, ErrorUnitialized
	}

	if teacher, ok := client.user.DBI.Teachers[id]; ok {
		return teacher, nil
	}
	return model.Teacher{}, ErrorNotFound
}

// GetClassroomByID is used to retrieve the classroom by it's specified ID.
// Returns ErrorNotFound if the classroom can't be found.
// Returns ErrorUnitialized if the user object hasn't been initialized.
func (client *EdupageClient) GetClassroomByID(id string) (model.Classroom, error) {
	if client.user == nil {
		return model.Classroom{}, ErrorUnitialized
	}

	if teacher, ok := client.user.DBI.Classrooms[id]; ok {
		return teacher, nil
	}
	return model.Classroom{}, ErrorNotFound
}

// FetchHomeworkAttachmens obtains the homework attchments for the specified homework.
// Returns ErrUnobtainableAttachments in case the attachments are not present.
// Retruns map, key is the resource name and value is the resource link
func (client *EdupageClient) FetchHomeworkAttachments(i *model.Homework) (map[string]string, error) {
	if len(i.ESuperID) == 0 || len(i.TestID) == 0 {
		return nil, errors.New("required fields superid and testid not set")
	}

	data := map[string]string{
		"testid":  i.TestID,
		"superid": i.ESuperID,
	}

	payload, err := CreatePayload(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload: %w", err)
	}

	resp, err := client.Credentials.httpClient.PostForm(
		"https://"+path.Join(client.Credentials.Server, "elearning", "?cmd=MaterialPlayer&akcia=getETestData"),
		payload,
	)
	if err != nil {
		return nil, fmt.Errorf("homework request failed: %w", err)
	}

	response, err := io.ReadAll(resp.Body)

	if len(response) < 5 {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	response = response[4:]

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(response)))
	_, err = base64.StdEncoding.Decode(decoded, response)
	if err != nil {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	decoded = bytes.Trim(decoded, "\x00")
	var object map[string]interface{}
	err = json.Unmarshal(decoded, &object)
	if err != nil {
		return nil, fmt.Errorf("homework request failed, bad response: %w", err)
	}

	attachments := make(map[string]string)

	// God help those who may try to debug this.
	if object["materialData"] == nil ||
		(reflect.TypeOf(object["materialData"]).Kind() != reflect.Map ||
			reflect.TypeOf(object["materialData"]).Elem().Kind() != reflect.Interface) {
		return nil, model.ErrUnobtainableAttachments
	}
	materialData := object["materialData"].(map[string]interface{})

	if materialData["cardsData"] == nil ||
		(reflect.TypeOf(materialData["cardsData"]).Kind() != reflect.Map ||
			reflect.TypeOf(materialData["cardsData"]).Elem().Kind() != reflect.Interface) {
		return nil, model.ErrUnobtainableAttachments
	}
	cardsData := materialData["cardsData"].(map[string]interface{})

	for _, entry := range cardsData {
		if entry == nil ||
			(reflect.TypeOf(entry).Kind() != reflect.Map ||
				reflect.TypeOf(entry).Elem().Kind() != reflect.Interface) {
			return nil, model.ErrUnobtainableAttachments
		}

		if e, ok := entry.(map[string]interface{})["content"]; !ok && reflect.TypeOf(e).Kind() != reflect.String {
			return nil, model.ErrUnobtainableAttachments
		}

		var content map[string]interface{}
		contentJson := entry.(map[string]interface{})["content"].(string)
		err = json.Unmarshal([]byte(contentJson), &content)
		if err != nil {
			return nil, err
		}

		if content["widgets"] == nil ||
			(reflect.TypeOf(content["widgets"]).Kind() != reflect.Slice ||
				reflect.TypeOf(content["widgets"]).Elem().Kind() != reflect.Interface) {
			return nil, model.ErrUnobtainableAttachments
		}

		widgets := content["widgets"].([]interface{})
		for _, widget := range widgets {
			if widget == nil ||
				(reflect.TypeOf(widget).Kind() != reflect.Map ||
					reflect.TypeOf(widget).Elem().Kind() != reflect.Interface) {
				return nil, model.ErrUnobtainableAttachments
			}
			if widget.(map[string]interface{})["props"] == nil ||
				(reflect.TypeOf(widget.(map[string]interface{})["props"]).Kind() != reflect.Map ||
					reflect.TypeOf(widget.(map[string]interface{})["props"]).Elem().Kind() != reflect.Interface) {
				return nil, model.ErrUnobtainableAttachments
			}
			props := widget.(map[string]interface{})["props"].(map[string]interface{})
			if files, ok := props["files"]; ok {
				for _, file := range files.([]interface{}) {
					if file == nil ||
						(reflect.TypeOf(file).Kind() != reflect.Map ||
							reflect.TypeOf(file).Elem().Kind() != reflect.Interface) {
						return nil, model.ErrUnobtainableAttachments
					}
					attachments[file.(map[string]interface{})["name"].(string)] = file.(map[string]interface{})["src"].(string)
				}
			}
		}
		if err != nil {
			continue
		}
		continue
	}

	return attachments, nil
}

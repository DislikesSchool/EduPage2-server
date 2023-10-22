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
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

var (
	ErrorUnitialized  = errors.New("unitialized")
	ErrorNotFound     = errors.New("not found")
	ErrorUnauthorized = errors.New("unauthorized")
	ErrorUnchangeable = errors.New("can not make changes at this time")
)

// EdupageClient is used to access the edupage api.
type EdupageClient struct {
	Credentials Credentials
	gsechash    string

	user    *model.User
	canteen *Canteen
	//timeline  *model.Timeline
	//results   *model.Results
	//timetable *model.Timetable
}

// CreateClient is used to create a client struct
func CreateClient(credentials Credentials) (*EdupageClient, error) {
	var client EdupageClient
	if credentials.httpClient == nil {
		return nil, errors.New("http client in credentials can not be nil")
	}
	client.Credentials = credentials

	user, err := client.fetchUserModel()
	if err != nil {
		return nil, err
	}

	client.user = &user

	return &client, nil
}

// UpdateCredentials updates the credentials and allows this struct to continue
// working after token expiry.
func (client *EdupageClient) UpdateCredentials(credentials Credentials) {
	client.Credentials = credentials
}

// GetUser retrieves the user from edupage or returns the stored data.
// If update is set to true, user data wil explicitly update
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetUser(update bool) (model.User, error) {
	if client.user == nil || update {
		user, err := client.fetchUserModel()
		if err != nil {
			return model.User{}, err
		}

		client.user = &user
		return user, nil
	} else {
		return *client.user, nil
	}
}

// GetRecentTimeline retrieves last 30 days of timeline from edupage.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetRecentTimeline() (model.Timeline, error) {
	timeline, err := client.fetchTimelineModel(time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return model.Timeline{}, err
	}

	return timeline, nil
}

// GetUser retrieves the timeline in a specified time interval from edupage.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetTimeline(from, to time.Time) (model.Timeline, error) {
	tt, err := client.fetchTimelineModel(from, to)
	if err != nil {
		return model.Timeline{}, err
	}
	return tt, nil
}

// GetRecentResults retrieves the results from the current year from edupage.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetRecentResults() (model.Results, error) {
	year := time.Now().Format("2006")
	halfyear := "RX" //TODO
	return client.fetchResultsModel(year, halfyear)
}

// GetResults retrieves the results in a specified interval from edupage.
// Halfyears types are: P1 (first halfyear), P2 (second halfyear), RX (whole year)
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetResults(year, halfyear string) (model.Results, error) {
	results, err := client.fetchResultsModel(year, halfyear)
	if err != nil {
		return model.Results{}, err
	}

	return results, nil
}

// GetResults retrieves this week's timetable from edupage.
func (client *EdupageClient) GetRecentTimetable() (model.Timetable, error) {
	tt, err := client.fetchTimetableModel(time.Now().AddDate(0, 0, -2), time.Now().AddDate(0, 0, 7))
	if err != nil {
		return model.Timetable{}, err
	}

	return tt, nil
}

// GetResults retrieves the timetable in the specified interval from edupage.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetTimetable(from, to time.Time) (model.Timetable, error) {
	tt, err := client.fetchTimetableModel(from, to)
	if err != nil {
		return model.Timetable{}, err
	}

	return tt, nil
}

// GetCanteen retrieves the whole week's canteen from the specified day.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetCanteen(date time.Time) (Canteen, error) {
	model, err := client.fetchCanteenModel(date)
	if err != nil {
		return Canteen{}, err
	}
	canteen, err := CreateCanteen(model)
	if err != nil {
		return Canteen{}, err
	}
	client.canteen = &canteen
	return canteen, nil
}

// GetCanteen retrieves the current week's canteen menu.
// Return ErrorUnauthorized if an authorization error occcurs.
func (client *EdupageClient) GetRecentCanteen() (Canteen, error) {
	day := time.Now().Weekday()
	if day == time.Saturday {
		return client.GetCanteen(time.Now().AddDate(0, 0, 2))
	} else if day == time.Sunday {
		return client.GetCanteen(time.Now().AddDate(0, 0, 1))
	}
	return client.GetCanteen(time.Now())
}

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
func (client *EdupageClient) FetchHomeworkAttachments(i model.Homework) (map[string]string, error) {
	if len(i.ESuperID) == 0 || len(i.TestID) == 0 {
		return nil, errors.New("required fields superid and testid not set")
	}

	data := map[string]string{
		"testid":  i.TestID,
		"superid": i.ESuperID,
	}

	payload := CreatePayload(data)

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

// ChangeOrderStatus changed order status of a meal for the specified day
// Return ErrorUnathorized, ErrorUnitialized, ErrorUnchangeable
func (e *EdupageClient) ChangeOrderStatus(day Day, order bool) error {
	if e.Credentials.httpClient == nil {
		return ErrorUnitialized
	}

	if !order && time.Now().After(day.CancelableUntil) {
		return ErrorUnchangeable
	}

	if order && time.Now().After(day.OrderableUntil) {
		return ErrorUnchangeable
	}

	var fids map[string]string
	var action string

	if order {
		fids = map[string]string{"2": "A"}
		action = "prihlas_do"
	} else {
		fids = map[string]string{"2": "AX"}
		action = "odhlas_do"
	}

	jedlaStravnika, _ := json.Marshal(CanteenPayload{
		BoarderID:   e.canteen.model.Info.BoarderID,
		BoarderUser: e.user.UserRow.UserID,
		Date:        day.Date.Format(model.TimeFormatYearMonthDay),
		FIDS:        fids, //TODO may be wrong
		View:        "pc_listok",
		Permission:  "Student",
		Action:      action,
	})

	payload := CreatePayload(map[string]string{
		"akcia":          "ulozJedlaStravnika",
		"jedlaStravnika": string(jedlaStravnika),
	})

	response, err := e.Credentials.httpClient.PostForm(fmt.Sprintf("https://%s/menu/", e.Credentials.Server), payload)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("invalid response code")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	decoded_body := make([]byte, base64.StdEncoding.DecodedLen(len(body)-4))

	_, err = base64.StdEncoding.Decode(decoded_body, body[4:])
	if err != nil {
		return fmt.Errorf("failed to decode response body: %s", err)
	}

	decoded_body = bytes.Trim(decoded_body, "\x00")

	var parsed map[string]interface{}
	err = json.Unmarshal(decoded_body, &parsed)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if parsed["status"] != nil {
		if reflect.TypeOf(parsed["status"]).Kind() != reflect.String {
			return errors.New("invalid response")
		}

		if parsed["status"].(string) == "insufficient_privileges" {
			return ErrorUnauthorized
		}
	}

	return nil
}

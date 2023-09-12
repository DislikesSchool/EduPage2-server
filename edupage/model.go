package edupage

import (
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

// This will contain datafixed and monolithic structures for ease of access

type Subject struct {
}

type Class struct {
}

type Teacher struct {
}

type Classroom struct {
}

type Student struct {
}

type TimetableRow struct {
	Type         string
	Date         time.Time
	Period       int
	StartTime    time.Time
	EndTime      time.Time
	Subject      Subject
	Classes      []Class
	Groups       []string
	IGroupID     string
	Teachers     []Teacher
	ClassroomIDs []Classroom
	Students     []Student
	Duration     int
}

func CreateTimetableRow(m model.Period) {
	//TODO
}

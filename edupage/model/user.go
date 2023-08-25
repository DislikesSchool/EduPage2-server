package model

type User struct {
	Edubar           map[string]interface{} `json:"_edubar"`
	Timeline         []TimelineItem         `json:"items"` // Only recent timeline, see EdupageClient.Timeline for more
	DBI              UserDataDBI            `json:"dbi"`
	UserRow          UserRow                `json:"userrow"`
	EventTypes       []EventType            `json:"eventtypes"`
	UserGroups       []string               `json:"usergroups"`
	NamedayToday     string                 `json:"meninyDnes"`
	NamedayTommorrow string                 `json:"meninyZajtra"`
}

type UserRow struct {
	UserID    string `json:"UserID"`
	StudentID string `json:"StudentID"`
	Firstname string `json:"p_meno"`
	Lastname  string `json:"p_priezvisko"`
	Email     string `json:"p_mail"`
	ClassID   string `json:"TriedaID"`
}

type EventType struct {
	ID            string   `json:"id"`
	C             string   `json:"c"`
	Name          string   `json:"name"`
	TTCancel      bool     `json:"ttcancel"`
	CTCan         bool     `json:"ctcan"`
	ClassRequired bool     `json:"classrequired"`
	Publish       string   `json:"publish"`
	NoCustomTime  bool     `json:"nocustomtime"`
	HideFields    []string `json:"hidefields"`
	DP            bool     `json:"dp"`
	Lesson        bool     `json:"lesson"`
	Attendance    bool     `json:"attendance"`
	DPrivacy      string   `json:"d_privacy"`
	CTEvent       bool     `json:"ctevent"`
	TemplateID    string   `json:"templateid"`
	CategoryID    string   `json:"categoryid"`
	SubID         string   `json:"subId"`
}

type UserDataDBI struct {
	Teachers           map[string]Teacher           `json:"teachers"`
	Classes            map[string]Class             `json:"classes"`
	Subjects           map[string]Subject           `json:"subjects"`
	Classrooms         map[string]Classrom          `json:"classrooms"`
	Students           map[string]Students          `json:"students"`
	Parents            map[string]Parents           `json:"parents"`
	Periods            []Period                     `json:"periods"`
	DayParts           map[string]DayParts          `json:"dayparts"`
	AbsentTypes        map[string]AbsentType        `json:"absenttypes"`
	SubstitutionTypes  map[string]SubstitionType    `json:"substitutiontypes"`
	StudentAbsentTypes map[string]StudentAbsentType `json:"studentabsenttypes"`
	EventTypes         map[string]UserEventType     `json:"eventtypes"`
	ProcessTypes       map[string]ProcessType       `json:"processtypes"`
	ProcessStates      map[string]ProcessState      `json:"processstates"`
	IsStudentAdult     bool                         `json:"isstudentadult"`
}

type Teacher struct {
	ID          string `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Short       string `json:"short"`
	Gender      string `json:"gender"`
	ClassroomID string `json:"classroomid"`
	DateFrom    string `json:"datefrom"`
	DateTo      string `json:"dateto"`
	IsOut       bool   `json:"isout"`
}

type Class struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Short       string `json:"short"`
	Grade       string `json:"grade"`
	TeacherID   string `json:"teacherid"`
	Teacher2ID  string `json:"teacher2id"`
	ClassroomID string `json:"classroomid"`
}

type Subject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Short    string `json:"short"`
	CBHidden bool   `json:"cbhidden"`
}

type Classrom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type Students struct {
	ID            string `json:"id"`
	ClassroomID   string `json:"classroomid"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	Parent1ID     string `json:"parent1id"`
	Parent2ID     string `json:"parent2id"`
	Parent3ID     string `json:"parent3id"`
	Gender        string `json:"gender"`
	DateFrom      string `json:"datefrom"`
	DateTo        string `json:"dateto"`
	NumberInClass string `json:"numberinclass"`
	IsOut         bool   `json:"isout"`
	Number        string `json:"number"`
	DataCopy      string `json:"kopiadata"`
}

type Parents struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
}

type Period struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type DayParts struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type AbsentType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type SubstitionType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type StudentAbsentType struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Short      string `json:"short"`
	Color      string `json:"color"`
	ExcuseType string `json:"excusetype"`
}

type UserEventType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type ProcessType struct {
	ID           string   `json:"id"`
	User         string   `json:"user"`
	Name         string   `json:"name"`
	Workflow     string   `json:"workflow"`
	Enabled      bool     `json:"enabled"`
	TextOptional bool     `json:"textoptional"`
	DataColumns  []string `json:"datacolumns"`
}

type ProcessState struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Icon    string            `json:"icon"`
	Color   string            `json:"color"`
	Next    map[string]string `json:"next"`
	Changes map[string]string `json:"changes"`
}

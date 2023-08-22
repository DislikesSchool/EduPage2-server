package edupage

type RawDataObject struct {
	Edubar     map[string]interface{} `json:"_edubar"`
	Items      []UserDataItem         `json:"items"`
	DBI        UserDataDBI            `json:"dbi"`
	UserRow    UserDataUserRow        `json:"userrow"`
	EventTypes []UserDataEvent        `json:"eventtypes"`
	UserGroups []string               `json:"usergroups"`
}

type UserDataItem struct {
	TimelineID      string `json:"timelineid"`
	Timestamp       string `json:"timestamp"`
	ReactionTo      string `json:"reakcia_na"`
	Typ             string `json:"typ"`
	User            string `json:"user"`
	TargetUser      string `json:"target_user"`
	UserName        string `json:"user_meno"`
	IneId           string `json:"ineid"`
	Text            string `json:"text"`
	AdditionTime    string `json:"cas_pridania"`
	EventTime       string `json:"cas_udalosti"`
	Data            string `json:"data"`
	Owner           string `json:"vlastnik"`
	OwnerName       string `json:"vlastnik_meno"`
	ReactionCount   string `json:"pocet_reakcii"`
	LastReaction    string `json:"posledna_reakcia"`
	HelpfulRecord   string `json:"pomocny_zaznam"`
	Removed         string `json:"removed"`
	AdditionTimeBtc string `json:"cas_pridania_btc"`
	LastReactionBtc string `json:"posledna_reakcia_btc"`
}

type UserDataUserRow struct {
	UserID    string `json:"UserID"`
	StudentID string `json:"StudentID"`
	Firstname string `json:"p_meno"`
	Lastname  string `json:"p_priezvisko"`
	Email     string `json:"p_mail"`
	ClassID   string `json:"TriedaID"`
}

type UserDataEvent struct {
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
	Teachers           map[string]UserDataTeacher           `json:"teachers"`
	Classes            map[string]UserDataClass             `json:"classes"`
	Subjects           map[string]UserDataSubject           `json:"subjects"`
	Classrooms         map[string]UserDataClassroom         `json:"classrooms"`
	Students           map[string]UserDataStudent           `json:"students"`
	Parents            map[string]UserDataParent            `json:"parents"`
	Periods            []UserDataPeriod                     `json:"periods"`
	DayParts           map[string]UserDataDayPart           `json:"dayparts"`
	AbsentTypes        map[string]UserDataAbsentType        `json:"absenttypes"`
	SubstitutionTypes  map[string]UserDataSubstitutionType  `json:"substitutiontypes"`
	StudentAbsentTypes map[string]UserDataStudentAbsentType `json:"studentabsenttypes"`
	EventTypes         map[string]UserDataEventType         `json:"eventtypes"`
	ProcessTypes       map[string]UserDataProcessType       `json:"processtypes"`
	ProcessStates      map[string]UserDataProcessState      `json:"processstates"`
	IsStudentAdult     bool                                 `json:"isstudentadult"`
}

type UserDataTeacher struct {
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

type UserDataClass struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Short       string `json:"short"`
	Grade       string `json:"grade"`
	TeacherID   string `json:"teacherid"`
	Teacher2ID  string `json:"teacher2id"`
	ClassroomID string `json:"classroomid"`
}

type UserDataSubject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Short    string `json:"short"`
	CBHidden bool   `json:"cbhidden"`
}

type UserDataClassroom struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type UserDataStudent struct {
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

type UserDataParent struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
}

type UserDataPeriod struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type UserDataDayPart struct {
	ID        string `json:"id"`
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
	Name      string `json:"name"`
	Short     string `json:"short"`
}

type UserDataAbsentType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type UserDataSubstitutionType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
}

type UserDataStudentAbsentType struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Short      string `json:"short"`
	Color      string `json:"color"`
	ExcuseType string `json:"excusetype"`
}

type UserDataEventType struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UserDataProcessType struct {
	ID           string   `json:"id"`
	User         string   `json:"user"`
	Name         string   `json:"name"`
	Workflow     string   `json:"workflow"`
	Enabled      bool     `json:"enabled"`
	TextOptional bool     `json:"textoptional"`
	DataColumns  []string `json:"datacolumns"`
}

type UserDataProcessState struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Icon    string            `json:"icon"`
	Color   string            `json:"color"`
	Next    map[string]string `json:"next"`
	Changes map[string]string `json:"changes"`
}

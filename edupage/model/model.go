package model

var (
	TimeFormatYearMonthDay = "2006-01-02"
)

type Mergeable interface {
	Merge(src *Mergeable)
}

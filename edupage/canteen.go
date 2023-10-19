package edupage

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage/model"
)

type Meal struct {
	Alergens []int
	Name     string
	Weight   int
}

type Menu struct {
	Meals []Meal
}

type Day struct {
	Date            time.Time
	AvailableFrom   time.Time
	AvailableTo     time.Time
	Ordered         bool
	OrderableUntil  time.Time
	CancelableUntil time.Time
	Menus           []Menu
}

// IsAvailable checks if the meal is currently available for consuming/pickup
func (m *Day) IsAvailable(t time.Time) bool {
	if t.After(m.AvailableFrom) && t.Before(m.AvailableTo) {
		return true
	} else {
		return false
	}
}

// CanOrder checks if current day's meal can still be ordered/unordered
func (m *Day) CanOrder(t time.Time) bool {
	if t.Before(m.OrderableUntil) {
		return false
	}

	if t.After(m.CancelableUntil) {
		return false
	}

	return true
}

func (m *Day) Cancel() {
	//TODO
}

func (m *Day) Order() {
	//TODO
}

// Represents the canteen, contains menu information, and additional information
type Canteen struct {
	Days map[string]Day
}

// Obtain the menu for a specified day.
// Returns menu, or false bool, indicating that no menu for that day was found.
func (c *Canteen) GetMenuByDay(time time.Time) (Day, bool) {
	if menu, exists := c.Days[time.Format("2006-01-02")]; exists {
		return menu, true
	} else {
		return Day{}, false
	}
}

// Global

// CreateDay creates a Day object from model.CanteenDay
func CreateDay(date string, day model.CanteenDay) (Day, error) {
	from, err := parseCanteenDate(date, day.AvailableFrom)
	if err != nil {
		return Day{}, err
	}

	to, err := parseCanteenDate(date, day.AvailableTo)
	if err != nil {
		return Day{}, err
	}

	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Day{}, err
	}

	orderable, err := time.Parse("2006-01-02 15:04", day.OrderableUntil)
	if err != nil {
		return Day{}, err
	}

	cancelable, err := time.Parse("2006-01-02 15:04", day.CancelableUntil)
	if err != nil {
		return Day{}, err
	}

	var meals []Meal = make([]Meal, len(day.Rows))
	for index, row := range day.Rows {
		alergens := make([]int, len(row.AlergenIDs))
		for k := range row.AlergenIDs {
			n, _ := k.Int64()
			alergens = append(alergens, int(n))
		}

		weight, _ := row.Weights.Int64()

		meals[index] = Meal{
			Alergens: alergens,
			Name:     row.Name,
			Weight:   int(weight),
		}
	}

	return Day{
		Date:            dt,
		AvailableFrom:   from,
		AvailableTo:     to,
		Ordered:         day.Evidence.Status == "A",
		OrderableUntil:  orderable,
		CancelableUntil: cancelable,
		Menus:           []Menu{{Meals: meals}},
	}, nil
}

// CreateCanteen creates Canteen object from model.Canteen
func CreateCanteen(m model.Canteen) (Canteen, error) {
	days := map[string]Day{}
	for date, day := range m.Days {
		day, err := CreateDay(date, day)
		if err == nil {
			days[date] = day
		} else {
			return Canteen{}, err
		}
	}

	return Canteen{
		Days: days,
	}, nil
}

// PRIVATE

func parseCanteenDate(date, hm string) (time.Time, error) {
	//TODO: maybe regex?
	split := strings.Split(hm, ":")

	hour, err := strconv.Atoi(split[0])
	if err != nil {
		return time.Time{}, errors.New("failed to parse hours")
	}

	minute, err := strconv.Atoi(split[1])
	if err != nil {
		return time.Time{}, errors.New("failed to parse minutes")
	}

	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(dt.Year(),
		dt.Month(),
		dt.Day(),
		hour,
		minute,
		0,
		0,
		time.Now().UTC().Location()), nil
}

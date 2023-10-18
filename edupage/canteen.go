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
	Date          time.Time
	AvailableFrom time.Time
	AvailableTo   time.Time
	Ordered       bool
	OrderFrom     time.Time
	OrderUntil    time.Time
	Meals         []Meal
}

func (m *Menu) IsAvailable(t time.Time) bool {
	if t.Before(m.AvailableFrom) {
		return false
	}

	if t.After(m.AvailableTo) {
		return false
	}

	return true
}

func (m *Menu) CanOrder(t time.Time) bool {
	if t.Before(m.OrderFrom) {
		return false
	}

	if t.After(m.OrderUntil) {
		return false
	}

	return true
}

func (m *Menu) Cancel() {

}

func (m *Menu) Order() {

}

func CreateMenu(date string, day model.CanteenDay) (Menu, error) {
	from, err := parseCanteenDate(date, day.AvailableFrom)
	if err != nil {
		return Menu{}, err
	}

	to, err := parseCanteenDate(date, day.AvailableTo)
	if err != nil {
		return Menu{}, err
	}

	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Menu{}, err
	}

	order_from, err := time.Parse("2006-01-02 15:04", day.OrderFrom)
	if err != nil {
		return Menu{}, err
	}

	order_until, err := time.Parse("2006-01-02 15:04", day.OrderUntil)
	if err != nil {
		return Menu{}, err
	}

	var meals []Meal = make([]Meal, len(day.Rows))
	for index, row := range day.Rows {
		alergens := make([]int, len(row.AlergenIDs))
		for k, _ := range row.AlergenIDs {
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

	return Menu{
		Date:          dt,
		AvailableFrom: from,
		AvailableTo:   to,
		Ordered:       day.Evidence.Status == "A",
		OrderFrom:     order_from,
		OrderUntil:    order_until,
		Meals:         meals,
	}, nil
}

func parseCanteenDate(date, hm string) (time.Time, error) {
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

// Represents the canteen, contains menu information, and additional information
type Canteen struct {
	Menus map[string]Menu
}

// Obtain the menu for a specified day.
// Returns menu, or false bool, indicating that no menu for that day was found.
func (c *Canteen) GetMenuByDay(time time.Time) (Menu, bool) {
	if menu, exists := c.Menus[time.Format("2006-01-02")]; exists {
		return menu, true
	} else {
		return Menu{}, false
	}
}

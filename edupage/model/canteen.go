package model

import (
	"encoding/json"
	"fmt"
)

type TicketRow struct {
	Histamin    int                  `json:"histamin"`
	ReceptureID json.Number          `json:"receptura_id"`
	AlergenIDs  map[json.Number]bool `json:"alegenyIDS"`
	Name        string               `json:"nazov"`
	Alergens    string               `json:"alergenyStr"`
	Weights     json.Number          `json:"hmotnostiStr"`
}

type MenuItem struct {
	Name          string               `json:"nazov"`
	MenuName      string               `json:"nazovMenu"`
	MenuShortName string               `json:"skratkaMenu"`
	Alergens      string               `json:"alergenyStr"`
	AlergenIDs    map[json.Number]bool `json:"alergenyIDS"`
	Weights       string               `json:"hmotnostiStr"`
	MeatOrigin    map[string]bool      `json:"povodMasa"`
}

type CanteenDay struct {
	AvailableFrom           string                   `json:"vydaj_od"`
	AvailableTo             string                   `json:"vydaj_do"`
	FoodTypes               json.Number              `json:"drugov_jedal"`
	IsCooking               bool                     `json:"isCooking"`
	IsPickUp                bool                     `json:"isPickUp"`
	IsChoosable             bool                     `json:"isChoosable"`
	VisibleMenus            map[json.Number]bool     `json:"visibleMenus"`
	ChoosableMenus          map[json.Number]bool     `json:"choosableMenus"`
	UnsubscribeOnlyWholeDay bool                     `json:"odhlasIbaCelyDen"`
	OrderFrom               string                   `json:"prihlas_od"`
	OrderUntil              string                   `json:"prihlas_do"`
	ChangeTo                string                   `json:"zmen_do"`
	Edupage                 string                   `json:"edupage"`
	IsSemiChoosable         string                   `json:"isSemiChoosable"`
	Rows                    []TicketRow              `json:"rows"`
	Name                    string                   `json:"nazov"`
	Alergens                string                   `json:"alergenyStr"`
	Menus                   map[json.Number]MenuItem `json:"menus"`
	MenuNames               map[json.Number]struct {
		Name      string `json:"nazov"`
		ShortName string `json:"skratka"`
	} `json:"nazvyMenu"`
	Evidence struct {
		Status       string `json:"stav"`
		Obj          string `json:"obj"`
		PovObj       string `json:"pov_obj"`
		RowChanged   string `json:"row_changed"`
		TZRowChanged string `json:"tz_row_changed"`
	} `json:"evidencia"`
	IsRating bool   `json:"isRating"`
	RateFrom string `json:"rate_od"`
	RateTo   string `json:"rate_do"`
}

type AllCreditInfo struct {
}

type BoarderRow struct {
	ID            string      `json:"stravnikid"`
	Edupage       string      `json:"edupage"`
	Type          string      `json:"typ"`
	AgendaID      string      `json:"agendaid"`
	Name          string      `json:"meno"`
	Surname       string      `json:"priezvisko"`
	Subname       interface{} `json:"subname"`
	Class         string      `json:"trieda"`
	Chips         string      `json:"cipy"`
	Credits       string      `json:"kredit"`
	Backup        interface{} `json:"zaloha"`
	AllCreditInfo string      `json:"allKreditInfo"`
	HardMode      bool        `json:"hardMode"`
	BoarderCount  json.Number `json:"pocetStravnikov"`
	Note          interface{} `json:"poznamka"`
	Alergens      interface{} `json:"alergeny"`
	Histamin      json.Number `json:"histamin"`
	Birthdate     string      `json:"birthdate"`
	DateTo        string      `json:"dateto"`
	AgeGroup      json.Number `json:"vekova_skupina"`
	Active        json.Number `json:"active"`
	Visible       json.Number `json:"visible"`
	DeletedByUser json.Number `json:"deletedByUser"`
	Status        string      `json:"stav"`
	RowChanged    string      `json:"row_changed"`
	User          string      `json:"user"`
	FullName      string      `json:"fullName"`
}

func (b *BoarderRow) GetAllCreditInfo() AllCreditInfo {
	var r AllCreditInfo
	UnmarshalNestedString(b.AllCreditInfo, &r)
	return r
}

type Info struct {
	Credit             float64     `json:"kredit"`
	CreditBackup       interface{} `json:"zalohaKredit"`
	FirstMinusDate     string      `json:"prvyMinusDate"`
	LastOkDate         string      `json:"poslednyOkDate"`
	InterruptDay       interface{} `json:"prerusDen"`
	HasInterruptObject interface{} `json:"maPreruseneObj"`
	BoarderID          string      `json:"stravnikid"`
	Info2              struct {
		Credit    float64       `json:"kredit"`
		PomCredit string        `json:"pomKredit"`
		DayCount  int           `json:"pocetDni"`
		DebugDays []interface{} `json:"debugDni"`
		Done      bool          `json:"done"`
	} `json:"poslednyOkDate"`
	HistaminShow string `json:"showOnlyAffected"`
	Alergens     map[json.Number]struct {
		ID      int    `json:"id"`
		Tag     string `json:"ozn"`
		Name    string `json:"nazov"`
		Visible bool   `json:"visible"`
	} `json:"alergenyIDS"`
	BoarderRow BoarderRow `json:"strRow"`
}

type Canteen struct {
	Info Info
	Days map[string]CanteenDay
}

func ParseCanteen(data []byte) (Canteen, error) {
	var response map[string]interface{}

	err := json.Unmarshal(data, &response)
	if err != nil {
		return Canteen{}, fmt.Errorf("failed to parse canteen data: %s", err.Error())
	}

	edupage := response["csssnina"].(map[string]interface{})
	novyListok := edupage["novyListok"].(map[string]interface{})
	var info Info
	days := make(map[string]CanteenDay, len(novyListok)-1)
	for k, v := range novyListok {
		if k == "addInfo" {
			b, _ := json.Marshal(v)
			json.Unmarshal(b, &info)
		} else {
			var day CanteenDay
			b, _ := json.Marshal(v.(map[string]interface{})["2"])
			json.Unmarshal(b, &day)
			days[k] = day
		}
	}

	return Canteen{}, nil
}

func keys[K comparable, V any](data map[K]V) []K {
	var keys = make([]K, len(data))

	for k, _ := range data {
		keys = append(keys, k)
	}

	return keys
}

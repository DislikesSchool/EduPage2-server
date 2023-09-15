package model

import "encoding/json"

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

type Ticket struct {
	AvailableFrom           string                   `json:"vydaj_od"`
	AvailableTo             string                   `json:"vydaj_do"`
	FoodTypes               json.Number              `json:"drugov_jedal"`
	IsCooking               bool                     `json:"isCooking"`
	IsPickUp                bool                     `json:"isPickUp"`
	IsChoosable             bool                     `json:"isChoosable"`
	VisibleMenus            map[json.Number]bool     `json:"visibleMenus"`
	ChoosableMenus          map[json.Number]bool     `json:"choosableMenus"`
	UnsubscribeOnlyWholeDay bool                     `json:"odhlasIbaCelyDen"`
	SubscribeFrom           string                   `json:"prihlas_od"`
	SubscribeTo             string                   `json:"prihlas_do"`
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

type Canteen struct {
	Tickets map[string]Ticket
}

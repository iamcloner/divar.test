package schema

type Filters struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

type Subs struct {
	Code     int16     `json:"code"`
	Name     string    `json:"name"`
	Enable   bool      `json:"enable"`
	Filters  []Filters `json:"filters"`
	Requires []string  `json:"requires"`
}

type Categories struct {
	Code   int16  `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
	Subs   []Subs `json:"subs"`
}

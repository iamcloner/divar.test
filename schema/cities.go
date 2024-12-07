package schema

type Zones struct {
	Code   int16  `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type Cities struct {
	Code   int16   `json:"code"`
	Name   string  `json:"name"`
	Enable bool    `json:"enable"`
	Zones  []Zones `json:"zones"`
}

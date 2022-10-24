package config

type Object struct {
	ID    ObjectID `json:"id"`
	Name  string   `json:"name,omitempty"`
	Title string   `json:"title,omitempty"`

	COV *COV `json:"COV,omitempty"`

	Priorities []Priority `json:"priorities,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

type Priority struct {
	Name  string `json:"name,omitempty"`
	Level uint8  `json:"level,omitempty"`
}

type Property struct {
	Name string     `json:"name,omitempty"`
	ID   PropertyID `json:"id,omitempty"`
}

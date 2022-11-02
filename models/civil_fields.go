package models

type Field struct {
	Type       string            `json:"type"`
	Label      string            `json:"label"`
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
}

type CivilFields struct {
	Group  string  `json:"group"`
	Fields []Field `json:"fields"`
}

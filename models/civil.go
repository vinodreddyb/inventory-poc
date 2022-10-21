package models

type Civil struct {
	Id       string                 `bson:"_id" json:"id,omitempty"`
	Name     string                 `json:"name"`
	Unit     string                 `json:"unit"`
	Quantity int                    `json:"quantity"`
	Path     string                 `json:"path"`
	Supply   map[string]interface{} `json:"supply"`
	Install  map[string]interface{} `json:"install"`
}

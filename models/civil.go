package models

import (
	"encoding/json"
	"time"
)

type ISODate struct {
	Format string
	time.Time
}

//UnmarshalJSON ISODate method
func (Date *ISODate) UnmarshalJSON(b []byte) error {

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	Date.Format = "2006-01-02"
	t, _ := time.Parse(Date.Format, s)
	Date.Time = t
	return nil
}

// MarshalJSON ISODate method
func (Date ISODate) MarshalJSON() ([]byte, error) {
	return json.Marshal(Date.Time.Format(Date.Format))
}

type Civil struct {
	Id        string                 `bson:"_id" json:"id,omitempty"`
	Name      string                 `json:"name"`
	Unit      string                 `json:"unit"`
	Quantity  int                    `json:"quantity"`
	StartDate time.Time              `bson:"startDate,omitempty"json:"startDate,omitempty"`
	EndDate   time.Time              `bson:"endDate,omitempty" json:"endDate,omitempty"`
	Path      string                 `json:"path"`
	Supply    map[string]interface{} `json:"supply"`
	Install   map[string]interface{} `json:"install"`
	Tender    bool                   `json:"tender"`
}

type CivilDTO struct {
	Id        string                 `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Tender    bool                   `json:"tender"`
	Unit      string                 `json:"unit"`
	Quantity  int                    `json:"quantity"`
	StartDate ISODate                `json:"startDate,omitempty"`
	EndDate   ISODate                `json:"endDate,omitempty"`
	Path      string                 `json:"path"`
	Supply    map[string]interface{} `json:"supply"`
	Install   map[string]interface{} `json:"install"`
}

package models

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Type      string                 `json:"type"`
}

type CivilProgessJoin struct {
	Id        string                 `bson:"_id" json:"id,omitempty"`
	Name      string                 `json:"name"`
	Unit      string                 `json:"unit"`
	Quantity  int                    `json:"quantity"`
	StartDate time.Time              `bson:"startDate,omitempty"json:"startDate,omitempty"`
	EndDate   time.Time              `bson:"endDate,omitempty" json:"endDate,omitempty"`
	Path      string                 `json:"path"`
	Supply    map[string]interface{} `json:"supply"`
	Install   map[string]interface{} `json:"install"`
	Type      string                 `json:"type"`
	Progress  []CivilProgress
}
type CivilProgressDTO struct {
	Id         string  `json:"id,omitempty"`
	NodeId     string  `json:"nodeId"`
	Date       ISODate `json:"date,omitempty"`
	Percentage float64 `json:"percentage"`
}

type CivilProgress struct {
	Id         *primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	NodeId     string
	Date       time.Time
	Percentage float64
}

type CivilDTO struct {
	Id        string                 `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Unit      string                 `json:"unit"`
	Quantity  int                    `json:"quantity"`
	StartDate ISODate                `json:"startDate,omitempty"`
	EndDate   ISODate                `json:"endDate,omitempty"`
	Path      string                 `json:"path"`
	Supply    map[string]interface{} `json:"supply"`
	Install   map[string]interface{} `json:"install"`
	Progress  []CivilProgressDTO     `json:"progress,omitempty"`
}

type CivilProgressGraph struct {
	Labels   []string       `json:"labels"`
	DataSets []GraphDataSet `json:"datasets"`
}

type GraphDataSet struct {
	Label           string    `json:"label"`
	BackgroundColor string    `json:"backgroundColor"`
	Data            []float64 `json:"data"`
}

func CivilProgressDtoToDo(civil CivilProgressDTO) CivilProgress {
	civilDto := CivilProgress{
		NodeId:     civil.NodeId,
		Date:       civil.Date.Time,
		Percentage: civil.Percentage,
	}
	return civilDto
}

func CivilProgressDoToDto(civil CivilProgress) CivilProgressDTO {
	civilDto := CivilProgressDTO{
		NodeId:     civil.NodeId,
		Date:       ISODate{Time: civil.Date, Format: "2006-01-02"},
		Percentage: civil.Percentage,
	}
	return civilDto
}
func CivilDoToDto(civil Civil) CivilDTO {
	civilDto := CivilDTO{
		Id:        civil.Id,
		Name:      civil.Name,
		Type:      civil.Type,
		Unit:      civil.Unit,
		Quantity:  civil.Quantity,
		StartDate: ISODate{Time: civil.StartDate, Format: "2006-01-02"},
		EndDate:   ISODate{Time: civil.EndDate, Format: "2006-01-02"},
		Path:      civil.Path,
		Supply:    civil.Supply,
		Install:   civil.Install,
	}
	return civilDto
}
func CivilJoinDoToDto(civil CivilProgessJoin) CivilDTO {
	var progressDTOS []CivilProgressDTO
	for _, pr := range civil.Progress {
		progressDTOS = append(progressDTOS, CivilProgressDoToDto(pr))
	}
	civilDto := CivilDTO{
		Id:        civil.Id,
		Name:      civil.Name,
		Type:      civil.Type,
		Unit:      civil.Unit,
		Quantity:  civil.Quantity,
		StartDate: ISODate{Time: civil.StartDate, Format: "2006-01-02"},
		EndDate:   ISODate{Time: civil.EndDate, Format: "2006-01-02"},
		Path:      civil.Path,
		Supply:    civil.Supply,
		Install:   civil.Install,
		Progress:  progressDTOS,
	}
	return civilDto
}
func CivilDtoToDo(civil CivilDTO) Civil {
	civilDto := Civil{
		Id:        civil.Id,
		Name:      civil.Name,
		Type:      civil.Type,
		Unit:      civil.Unit,
		Quantity:  civil.Quantity,
		StartDate: civil.StartDate.Time,
		EndDate:   civil.EndDate.Time,
		Path:      civil.Path,
		Supply:    civil.Supply,
		Install:   civil.Install,
	}
	return civilDto
}

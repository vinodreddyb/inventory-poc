package models

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContractActivities struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Name         string             `json:"name"`
	Weightage    float64            `json:"weightage"`
	UnitWeigtage float64            `json:"unitWeightage"`
	Unit         string             `json:"unit"`
	SORQuantity  int                `json:"sorQuantity"`
	SORCost      int                `json:"sorCost"`
	StartDate    time.Time          `bson:"startDate,omitempty" json:"startDate,omitempty"`
	FinishDate   time.Time          `bson:"finishDate,omitempty" json:"finishDate,omitempty"`
}

type ContractSchedule struct {
	ActivityId string  `json:"name"`
	Year       int     `json:"year"`
	Month      string  `json:"month"`
	Type       string  `json:"sch"`
	Value      float64 `json:"act"`
}

type ScheduleGraph struct {
	Labels        []string  `json:"labels"`
	SchProgress   []float64 `json:"schProgress"`
	MonthActual   []float64 `json:"monthlyActual"`
	SchCumulative []float64 `json:"schCumulative"`
	ActCumulative []float64 `json:"actualCumulative"`
}

type ContractProgress struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Name          string             `json:"name"`
	Weightage     float64            `json:"weightage"`
	UnitWeightage float64            `json:"unitWeightage"`
	Unit          string             `json:"unit"`
	SORQuantity   int                `json:"sorQuantity"`
	SORCost       int                `json:"sorCost"`
	StartDate     time.Time          `bson:"startDate,omitempty" json:"startDate,omitempty"`
	FinishDate    time.Time          `bson:"finishDate,omitempty" json:"finishDate,omitempty"`
	Schedules     []ContractSchedule
}

type ContractProgressDTO struct {
	Weightage     float64   `json:"weightage"`
	UnitWeightage float64   `json:"unitWeightage"`
	Unit          string    `json:"unit"`
	SORQuantity   int       `json:"sorQuantity"`
	SORCost       int       `json:"sorCost"`
	StartDate     time.Time `bson:"startDate,omitempty" json:"startDate,omitempty"`
	FinishDate    time.Time `bson:"finishDate,omitempty" json:"finishDate,omitempty"`
	Type          string    `json:"type"`
}

type ContractProgressResponse struct {
	Columns    []string
	SchColumns map[int][]string
	Data       []map[string]interface{}
}
type Schedule struct {
	Year  int
	Month []Months
}

type Months struct {
	Name  string
	Value float64
	Order int
}

func CreateStruct(fields []string, types []interface{}) reflect.Value {
	f := []reflect.StructField{}
	/*fields := []string{"Myfield", "Area", "Size"}
	types := []interface{}{25, 5, 10.3}*/

	for i, v := range fields {
		x := reflect.StructField{
			Name: reflect.ValueOf(v).Interface().(string),
			Type: reflect.TypeOf(types[i]),
		}

		f = append(f, x)
	}

	t := reflect.StructOf(f)

	e := reflect.New(t).Elem()

	return e
}

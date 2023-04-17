package models

import (
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
	ActivityId string `json:"name"`
	Year       int    `json:"year"`
	Month      string `json:"month"`
	Sch        int    `json:"sch"`
	Act        int    `json:"act"`
}

type ScheduleGraph struct {
	Labels        []string  `json:"labels"`
	SchProgress   []float64 `json:"schProgress"`
	MonthActual   []float64 `json:"monthlyActual"`
	SchCumulative []float64 `json:"schCumulative"`
	ActCumulative []float64 `json:"actualCumulative"`
}

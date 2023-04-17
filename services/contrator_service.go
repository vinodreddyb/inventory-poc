package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"mongo-rest/configs"
	"mongo-rest/models"
	"time"
)

var activityCollection = configs.GetCollection(configs.DB, "Activities1")
var activitySchduleCollection = configs.GetCollection(configs.DB, "ActivitySchedule2")

func CaluclateProgress() models.ScheduleGraph {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m := getActivities(ctx)

	//unitWeightage
	group := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"month", "$month"},
						{"year", "$year"},
						{"type", "$type"},
					},
				},
				{"data", bson.D{
					{"$addToSet", bson.D{
						{"act", "$activityId"},
						{"month", "$month"},
						{"type", "$type"},
						{"value", "$value"},
					}}}},
			},
		},
	}
	aggregate, err := activitySchduleCollection.Aggregate(ctx, mongo.Pipeline{group})
	if err != nil {
		fmt.Println(err)
	}
	defer aggregate.Close(ctx)

	var progress map[string]float64 = make(map[string]float64)
	var progressAct map[string]float64 = make(map[string]float64)

	var monthYear []int32
	for aggregate.Next(ctx) {
		var doc bson.M
		err := aggregate.Decode(&doc)
		if err != nil {
			panic(err)
		}
		year := doc["_id"].(primitive.M)["year"].(int32)
		month := doc["_id"].(primitive.M)["month"].(string)
		valType := doc["_id"].(primitive.M)["type"].(string)
		if year == 2023 && month == "MAR" {
			fmt.Println("Hi")
		}
		b := doc["data"]
		var sumSch float64
		var sumAct float64
		for _, v := range b.(primitive.A) { // use type assertion to loop over []interface{}

			actvityWeigtage := m[v.(primitive.M)["act"].(string)]

			val := v.(primitive.M)["value"].(float64)

			if valType == "Sch." {
				sumSch = sumSch + (actvityWeigtage * val)
			} else {
				sumAct = sumAct + (actvityWeigtage * val)
			}

		}
		key := fmt.Sprintf("%s-%d", month, year)

		if monthYear == nil {
			monthYear = append(monthYear, year)
		}
		if !slices.Contains(monthYear, year) {
			monthYear = append(monthYear, year)
		}
		if valType == "Sch." {
			progress[key] = sumSch
		} else {
			progressAct[key] = sumAct
		}

	}
	slices.Sort(monthYear)

	var cummlativeProgress float64 = 0
	var cummlativeActual float64 = 0

	graph := buildScurveGraphData(monthYear, progress, cummlativeProgress, progressAct, cummlativeActual)

	return graph
}

func buildScurveGraphData(monthYear []int32, progress map[string]float64, cummlativeProgress float64, progressAct map[string]float64, cummlativeActual float64) models.ScheduleGraph {
	var labels []string
	var schMonthly []float64
	var actMonthly []float64
	var schCumilative []float64
	var actCumulative []float64
	months := []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUNE", "JULY", "AUG", "SEPT", "OCT", "NOV", "DEC"}

	for _, y := range monthYear {
		for _, month := range months {
			key := fmt.Sprintf("%s-%d", month, y)

			if v, exists := progress[key]; exists {
				labels = append(labels, key)
				cummlativeProgress += v
				schCumilative = append(schCumilative, cummlativeProgress)
				schMonthly = append(schMonthly, v)
			}
			if v, exists := progressAct[key]; exists {

				if !slices.Contains(labels, key) {
					labels = append(labels, key)
				}
				cummlativeActual += v
				actCumulative = append(actCumulative, cummlativeActual)
				actMonthly = append(actMonthly, v)
			}

		}
	}

	fmt.Println(progress)
	fmt.Println(progressAct)

	graph := models.ScheduleGraph{
		Labels:        labels,
		SchProgress:   schMonthly,
		MonthActual:   actMonthly,
		SchCumulative: schCumilative,
		ActCumulative: actCumulative,
	}
	return graph
}

func getMonthNumber(month string) int {

	switch month {
	case "JAN":
		return 1
	case "FEB":
		return 2
	case "MAR":
		return 3
	case "APR":
		return 4
	case "MAY":
		return 5
	case "JUNE":
		return 6
	case "JULY":
		return 7
	case "AUG":
		return 8
	case "SEPT":
		return 9
	case "OCT":
		return 10
	case "NOV":
		return 11
	case "DEC":
		return 12
	default:
		return 0

	}
}

func getMonth(month int) string {

	switch month {
	case 1:
		return "JAN"
	case 2:
		return "FEB"
	case 3:
		return "MAR"
	case 4:
		return "APR"
	case 5:
		return "MAY"
	case 6:
		return "JUN"
	case 7:
		return "JUL"
	case 8:
		return "AUG"
	case 9:
		return "SEPT"
	case 10:
		return "OCT"
	case 11:
		return "NOV"
	case 12:
		return "DEC"
	default:
		return ""

	}
}
func getActivities(ctx context.Context) map[string]float64 {
	cur, err2 := activityCollection.Find(ctx, bson.M{})
	defer cur.Close(ctx)
	if err2 != nil {
		fmt.Println(err2)
		panic(err2)
	}

	var m map[string]float64 = make(map[string]float64)

	for cur.Next(ctx) {
		var doc bson.M
		err2 := cur.Decode(&doc)
		if err2 != nil {

		}
		m[doc["_id"].(primitive.ObjectID).Hex()] = doc["unitWeightage"].(float64)
	}
	return m
}

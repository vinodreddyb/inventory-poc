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
	"reflect"
	"strconv"
	"time"
)

var activityCollection = configs.GetCollection(configs.DB, "Activities1")
var activitySchduleCollection = configs.GetCollection(configs.DB, "ActivitySchedule2")

func GetContractProgress() models.ContractProgressResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addFields := bson.D{
		{
			"$addFields", bson.D{{"actId", bson.D{{"$toString", "$_id"}}}},
		},
	}
	lookup := bson.D{
		{
			"$lookup", bson.D{
				{"from", "ActivitySchedule2"},
				{"localField", "actId"},
				{"foreignField", "activityId"},
				{"as", "schedules"},
			},
		},
	}
	aggregate, err := activityCollection.Aggregate(ctx, mongo.Pipeline{addFields, lookup})
	if err != nil {
		panic(err)
	}

	var progressDtos []map[string]interface{}
	schColumns := make(map[int][]string)
	for aggregate.Next(ctx) {
		var progress models.ContractProgress
		err := aggregate.Decode(&progress)
		if err != nil {
			panic(err)
		}
		group := make(map[int]map[string]float64)
		actualGroup := make(map[int]map[string]float64)

		for _, sch := range progress.Schedules {
			if sch.Type == "Sch." {

				if val, ok := group[sch.Year]; ok {
					val[sch.Month] = sch.Value
				} else {
					group[sch.Year] = map[string]float64{
						sch.Month: sch.Value,
					}
				}

			} else if sch.Type == "Act." {
				if val, ok := actualGroup[sch.Year]; ok {
					val[sch.Month] = sch.Value
				} else {
					actualGroup[sch.Year] = map[string]float64{
						sch.Month: sch.Value,
					}
				}

			}

			if v, ok := schColumns[sch.Year]; ok {
				if !slices.Contains(v, sch.Month) {
					v = append(v, sch.Month)
					schColumns[sch.Year] = v
				}

			} else {
				schColumns[sch.Year] = append(schColumns[sch.Year], sch.Month)
			}

		}

		/*dtoSch := models.ContractProgressDTO{
			Name:          progress.Name,
			Weightage:     progress.Weightage,
			UnitWeightage: progress.UnitWeightage,
			Unit:          progress.Unit,
			SORQuantity:   progress.SORQuantity,
			SORCost:       progress.SORCost,
			StartDate:     progress.StartDate,
			FinishDate:    progress.FinishDate,
			Type:          "Sch",
			Schedules:     group,
		}*/
		dtoSch1 := map[string]interface{}{
			"Activity":      progress.Name,
			"Weightage":     progress.Weightage,
			"UnitWeightage": progress.UnitWeightage,
			"Unit":          progress.Unit,
			"SORQuantity":   progress.SORQuantity,
			"SORCost":       progress.SORCost,
			"StartDate":     models.ISODate{Time: progress.StartDate, Format: "02/Jan/2006"},
			"FinishDate":    models.ISODate{Time: progress.FinishDate, Format: "02/Jan/2006"},
			"Type":          "Sch",
		}
		for y, mon := range group {
			for k, m := range mon {
				key := fmt.Sprintf("%s-%d", k, y)
				dtoSch1[key] = m
			}
		}
		progressDtos = append(progressDtos, dtoSch1)
		dtoAct := map[string]interface{}{

			"Activity":      progress.Name,
			"Weightage":     0,
			"UnitWeightage": progress.UnitWeightage,
			"Unit":          "",
			"SORQuantity":   0,
			"SORCost":       0,
			"Type":          "Act",
		}
		for y, mon := range actualGroup {
			for k, m := range mon {
				key := fmt.Sprintf("%s-%d", k, y)
				dtoAct[key] = m
			}
		}
		progressDtos = append(progressDtos, dtoAct)
	}

	fields := getScheduleProgressDisplayFields()

	response := models.ContractProgressResponse{
		Columns:    fields,
		SchColumns: schColumns,
		Data:       progressDtos,
	}
	return response

}

func getScheduleProgressDisplayFields() []string {
	t := reflect.TypeOf(models.ContractProgressDTO{})

	names := make([]string, t.NumField())
	for i := range names {
		names[i] = t.Field(i).Name
	}
	return names
}

func getMonthlySchedule(years []int, group map[int][]models.Months) ([]string, []interface{}) {
	var yearFields []string
	var yearValues []interface{}
	for _, year := range years {
		var fields []string
		var fieldVals []interface{}

		for _, m := range group[year] {
			fields = append(fields, m.Name)
			fieldVals = append(fieldVals, m.Value)
		}
		monthStruct := models.CreateStruct(fields, fieldVals)
		yearFields = append(yearFields, strconv.Itoa(year))
		yearValues = append(yearValues, monthStruct)

	}
	return yearFields, yearValues
}
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

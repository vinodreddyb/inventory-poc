package services

import (
	"context"
	"errors"
	"fmt"
	logr "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"inventory-poc/configs"
	"inventory-poc/models"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
var userCollection = configs.GetCollection("Users")
var civilCollection = configs.GetCollection("civil")
var civilFieldsCollection = configs.GetCollection("civil_fields")
var civilProgressCollection = configs.GetCollection("civil_progress")
*/
func AddNewUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := configs.GetCollection("Users").InsertOne(ctx, newUser)

	if err != nil {
		return nil, err
	}
	newUser.Id = result.InsertedID.(primitive.ObjectID)

	return &newUser, nil
}

func GetAllUsers() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := configs.GetCollection("Users").Find(ctx, bson.M{})
	if err != nil {
		logr.Error(err)
	}

	defer cursor.Close(ctx)

	var users []models.User

	for cursor.Next(ctx) {
		var usr models.User
		if err = cursor.Decode(&usr); err != nil {
			log.Fatal(err)
		}
		users = append(users, usr)
	}
	return users
}

func GetCivils(path string) ([]models.CivilDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{}
	if path != "" {
		//m = bson.M{"path": path}
		match := bson.D{{"$match", bson.M{"path": path}}}
		pipeline = append(pipeline, match)
	}
	lookup := bson.D{{"$lookup", bson.M{"from": "civil_progress",
		"localField":   "_id",
		"foreignField": "nodeid",
		"as":           "progress"}}}
	pipeline = append(pipeline, lookup)

	cursor, err := configs.GetCollection("civil").Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var civils []models.CivilDTO

	for cursor.Next(ctx) {
		var civil models.CivilProgessJoin
		err = cursor.Decode(&civil)
		if err != nil {
			logr.Error(err)
			return nil, err
		}
		civils = append(civils, models.CivilJoinDoToDto(civil))
	}
	return civils, nil
}

func GetCivilFields() ([]models.CivilFields, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := configs.GetCollection("civil_fields").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(ctx)
	var civilFields []models.CivilFields

	for cursor.Next(ctx) {
		var civilField models.CivilFields
		if err = cursor.Decode(&civilField); err != nil {
			logr.Error(err)
		}
		civilFields = append(civilFields, civilField)
	}
	return civilFields, nil
}

func AddCivilNode(nodePath string, civilNode models.CivilDTO) ([]models.CivilDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//civilCollection.U
	defer cancel()
	cur, err2 := configs.GetCollection("civil").Find(ctx, bson.M{"path": nodePath})
	var civils []models.CivilDTO
	if err2 != nil {
		logr.Error(err2)
		return nil, err2
	}
	for cur.Next(ctx) {
		var civil models.Civil
		if err2 = cur.Decode(&civil); err2 != nil {
			logr.Error(err2)
		}
		civils = append(civils, models.CivilDoToDto(civil))
	}

	if len(civils) > 0 {
		lastChild := civils[len(civils)-1]
		split := strings.Split(lastChild.Id, "-")
		lastChildId, _ := strconv.Atoi(split[len(split)-1])
		split[len(split)-1] = strconv.Itoa(lastChildId + 1)

		fmt.Printf("last child %d", strings.Join(split, "-"))

		civilNode.Id = strings.Join(split, "-")

	} else {
		split := strings.Split(nodePath, ",")
		civilNode.Id = split[len(split)-2] + "-" + strconv.Itoa(1)
	}
	civilNode.Path = nodePath

	_, err := configs.GetCollection("civil").InsertOne(ctx, models.CivilDtoToDo(civilNode))

	if err != nil {
		return nil, err
	}
	civils = append(civils, civilNode)

	return civils, err
}

func UpdateCivilNode(civilNode models.CivilDTO) (models.CivilDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"name":      civilNode.Name,
			"unit":      civilNode.Unit,
			"quantity":  civilNode.Quantity,
			"supply":    civilNode.Supply,
			"install":   civilNode.Install,
			"startDate": civilNode.StartDate.Time,
			"endDate":   civilNode.EndDate.Time,
		},
	}
	result := configs.GetCollection("civil").FindOneAndUpdate(ctx, bson.M{"_id": civilNode.Id}, update)

	if result.Err() != nil {
		return models.CivilDTO{}, result.Err()
	}

	return civilNode, nil
}

func AddWorkStatus(progress models.CivilProgressDTO) (*models.CivilProgressDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nodeDetails := configs.GetCollection("civil").FindOne(ctx, bson.M{"_id": progress.NodeId})
	if nodeDetails != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if nodeDetails.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("No documents found with node " + progress.NodeId)
		}
		var node models.Civil
		errD := nodeDetails.Decode(&node)
		if errD != nil {
			return nil, errors.New("Error decoding node from DB " + progress.NodeId)
		}

		if !(progress.Date.Time.Before(node.EndDate) &&
			(progress.Date.Time.Equal(node.StartDate) || progress.Date.Time.After(node.StartDate))) {
			return nil, errors.New("enter status date with in start date and end date range")
		}

	}
	progressDO := models.CivilProgressDtoToDo(progress)
	filter := bson.M{"nodeid": progressDO.NodeId, "date": progressDO.Date}
	update := bson.D{{"$set", progressDO}}
	opts := options.Update().SetUpsert(true)
	result, err := configs.GetCollection("civil_progress").UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return nil, err
	}
	if result.UpsertedID != nil {
		progress.Id = result.UpsertedID.(primitive.ObjectID).String()
	}

	return &progress, nil
}

func GetStatusGraph(nodeId string) ([]models.CivilProgressDTO, *models.CivilProgressGraph, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nodeDetails := configs.GetCollection("civil").FindOne(ctx, bson.M{"_id": nodeId})
	if nodeDetails != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if nodeDetails.Err() == mongo.ErrNoDocuments {
			return nil, nil, errors.New("No documents found with node " + nodeId)
		}
	}

	cursor, errD := configs.GetCollection("civil_progress").Find(ctx, bson.M{"nodeid": nodeId})
	if errD != nil {
		return nil, nil, errD
	}

	defer cursor.Close(ctx)
	statusesMap := make(map[string]models.CivilProgressDTO)
	for cursor.Next(ctx) {
		var status models.CivilProgress
		if errD = cursor.Decode(&status); errD != nil {
			logr.Error(errD)
			return nil, nil, errD
		}
		statusesMap[status.Date.Format("2006-01-02")] = models.CivilProgressDoToDto(status)
	}

	var node models.Civil
	errDD := nodeDetails.Decode(&node)
	if errDD != nil {
		return nil, nil, errors.New("Error decoding node from DB " + nodeId)
	}
	var civilProgress []models.CivilProgressDTO
	var labels []string
	var data []float64
	for d := node.StartDate; d.After(node.EndDate) == false; d = d.AddDate(0, 0, 1) {
		labels = append(labels, d.Format("2006-01-02"))
		if val, ok := statusesMap[d.Format("2006-01-02")]; ok {
			civilProgress = append(civilProgress, val)
			data = append(data, val.Percentage)
		} else {
			civilProgress = append(civilProgress, models.CivilProgressDTO{
				NodeId:     nodeId,
				Date:       models.ISODate{Time: d, Format: "2006-01-02"},
				Percentage: 0,
			})
			data = append(data, 0)
		}
	}
	graph := models.CivilProgressGraph{
		Labels: labels,
		DataSets: []models.GraphDataSet{
			{
				Label:           "Status",
				BackgroundColor: "#42A5F5",
				Data:            data,
			},
		},
	}
	return civilProgress, &graph, nil
}

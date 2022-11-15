package services

import (
	"context"
	"errors"
	logr "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mongo-rest/configs"
	"mongo-rest/models"
	"time"
)

var userCollection = configs.GetCollection(configs.DB, "Users")
var civilCollection = configs.GetCollection(configs.DB, "civil")
var civilFieldsCollection = configs.GetCollection(configs.DB, "civil_fields")
var civilProgressCollection = configs.GetCollection(configs.DB, "civil_progress")

func AddNewUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return nil, err
	}
	newUser.Id = result.InsertedID.(primitive.ObjectID)

	return &newUser, nil
}

func GetAllUsers() []models.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := userCollection.Find(ctx, bson.M{})
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

func GetCivils(path string) (*[]models.CivilDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m := bson.M{}
	if path != "" {
		m = bson.M{"path": path}
	}
	cursor, err := civilCollection.Find(ctx, m)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var civils []models.CivilDTO

	for cursor.Next(ctx) {
		var civil models.Civil
		err = cursor.Decode(&civil)
		if err != nil {
			logr.Error(err)
			return nil, err
		}
		civils = append(civils, models.CivilDoToDto(civil))
	}
	return &civils, nil
}

func GetCivilFields() ([]models.CivilFields, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := civilFieldsCollection.Find(ctx, bson.M{})
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

func AddCivilNode(civilNode models.CivilDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//civilCollection.U
	defer cancel()
	_, err := civilCollection.InsertOne(ctx, models.CivilDtoToDo(civilNode))

	if err != nil {
		return err
	}

	return nil
}

func UpdateCivilNode(civilNode models.CivilDTO) (models.CivilDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{"unit": civilNode.Unit,
			"quantity":  civilNode.Quantity,
			"supply":    civilNode.Supply,
			"install":   civilNode.Install,
			"startDate": civilNode.StartDate.Time,
			"endDate":   civilNode.EndDate.Time,
		},
	}
	result := civilCollection.FindOneAndUpdate(ctx, bson.M{"_id": civilNode.Id}, update)

	if result.Err() != nil {
		return models.CivilDTO{}, result.Err()
	}

	return civilNode, nil
}

func AddWorkStatus(progress models.CivilProgressDTO) (*models.CivilProgressDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nodeDetails := civilCollection.FindOne(ctx, bson.M{"_id": progress.NodeId})
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

		if !(progress.Date.Time.Before(node.EndDate) && progress.Date.Time.After(node.StartDate)) {
			return nil, errors.New("enter status date with in start date and end date range")
		}

	}
	progressDO := models.CivilProgressDtoToDo(progress)
	result, err := civilProgressCollection.InsertOne(ctx, progressDO)

	if err != nil {
		return nil, err
	}
	progress.Id = result.InsertedID.(primitive.ObjectID).String()

	return &progress, nil
}

package services

import (
	"context"
	logr "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mongo-rest/configs"
	"mongo-rest/models"
	"time"
)

var userCollection = configs.GetCollection(configs.DB, "Users")
var civilCollection = configs.GetCollection(configs.DB, "civil")
var civilFieldsCollection = configs.GetCollection(configs.DB, "civil_fields")

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
		civils = append(civils, civilDoToDto(civil))
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
	_, err := civilCollection.InsertOne(ctx, civilDtoToDo(civilNode))

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

func civilDoToDto(civil models.Civil) models.CivilDTO {
	civilDto := models.CivilDTO{
		Id:        civil.Id,
		Name:      civil.Name,
		Unit:      civil.Unit,
		Quantity:  civil.Quantity,
		StartDate: models.ISODate{Time: civil.StartDate, Format: "2006-01-02"},
		EndDate:   models.ISODate{Time: civil.EndDate, Format: "2006-01-02"},
		Path:      civil.Path,
		Supply:    civil.Supply,
		Install:   civil.Install,
	}
	return civilDto
}
func civilDtoToDo(civil models.CivilDTO) models.Civil {
	civilDto := models.Civil{
		Id:        civil.Id,
		Name:      civil.Name,
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

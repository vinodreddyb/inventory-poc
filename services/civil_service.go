package services

import (
	"context"
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
		log.Fatal(err)
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

func GetCivils(path string) ([]models.Civil, error) {
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

	var civils []models.Civil

	for cursor.Next(ctx) {
		var civil models.Civil
		if err = cursor.Decode(&civil); err != nil {
			log.Fatal(err)
		}
		civils = append(civils, civil)
	}
	return civils, nil
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
			log.Fatal(err)
		}
		civilFields = append(civilFields, civilField)
	}
	return civilFields, nil
}

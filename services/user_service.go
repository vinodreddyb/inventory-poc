package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mongo-rest/configs"
	"mongo-rest/dbmodels"
	"mongo-rest/models"
	"time"
)

var userCollection = configs.GetCollection(configs.DB, "Users")
var civilCollection = configs.GetCollection(configs.DB, "civil")

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

func GetCivils() []dbmodels.Civil {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := civilCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var civils []dbmodels.Civil

	for cursor.Next(ctx) {
		var civil dbmodels.Civil
		if err = cursor.Decode(&civil); err != nil {
			log.Fatal(err)
		}
		civils = append(civils, civil)
	}
	return civils
}

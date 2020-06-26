package data

import (
	"context"
	"log"
	"time"

	"github.com/mauriliommachado/go-commerce/user-service/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var app *models.App

//InitDb initate the middleware with app configs
func InitDb(config *models.App) {
	app = config
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	app.Collection = client.Database("go-commerce").Collection(app.ServiceName)
}

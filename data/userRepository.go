package data

import (
	"context"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/mauriliommachado/go-commerce/user-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//UserRepository type
type UserRepository struct {
	C *mongo.Collection
}

//Create user
func (r *UserRepository) Create(user *models.User) error {
	insertResult, err := r.C.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	user.ID = insertResult.InsertedID.(primitive.ObjectID)
	log.Println("Inserted a single document: ", insertResult.InsertedID.(primitive.ObjectID).Hex())
	return err
}

//Update user
func (r *UserRepository) Update(user *models.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"username": user.Username,
		"password": user.Password}}
	_, err := r.C.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Update a single document ")
	return err
}

//Get user
func (r *UserRepository) Get(id string) (models.User, error) {
	// create a value into which the result can be decoded
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	var result models.User
	err := r.C.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
	}
	return result, err
}

//GetAll users
func (r *UserRepository) GetAll() []models.User {
	// Pass these options to the Find method
	findOptions := options.Find()

	// Here's an array in which you can store the decoded documents
	var results []models.User

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := r.C.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	return results
}

//Delete User
func (r *UserRepository) Delete(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	_, err := r.C.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

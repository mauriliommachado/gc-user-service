package models

import (
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

// App application manager
type App struct {
	Router      *httprouter.Router
	Collection  *mongo.Collection
	ServiceName string
}

package models

import (
	"go-react/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//EventType model
type EventType struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

//EventTypeQuery var is to query data from eventTypes collection
var EventTypeQuery = database.MongoQuery{CollectionName: "eventTypes"}

package models

import (
	"go-react/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Event model
type Event struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                  string             `json:"name,omitempty" bson:"name,omitempty"`
	Language              []string           `json:"language,omitempty" bson:"language,omitempty"`
	TypeID                primitive.ObjectID `json:"typeID,omitempty" bson:"typeID,omitempty"`
	Location              string             `json:"location,omitempty" bson:"location,omitempty"`
	Accommodation         string             `json:"accommodation,omitempty" bson:"accommodation,omitempty"`
	RegistrationCloseDate time.Time          `json:"registrationCloseDate,omitempty" bson:"registrationCloseDate,omitempty"`
	StartDate             time.Time          `json:"startDate,omitempty" bson:"startDate,omitempty"`
	EndDate               time.Time          `json:"endDate,omitempty" bson:"endDate,omitempty"`
	MaxParticipants       string             `json:"maxParticipants,omitempty" bson:"maxParticipants,omitempty"`
	Tags                  []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Description           string             `json:"description,omitempty" bson:"description,omitempty"`
	OwnerID               primitive.ObjectID `json:"ownerID,omitempty" bson:"ownerID,omitempty"`
	Budget                string             `json:"budget,omitempty" bson:"budget,omitempty"`
	IsApproved            bool               `json:"isApproved,omitempty" bson:"isApproved,omitempty"`
	Image                 string             `json:"image,omitempty" bson:"image,omitempty"`
	ReviewerID            primitive.ObjectID `json:"reviewerID,omitempty" bson:"reviewerID,omitempty"`
	IsFinished            bool               `json:"isFinished,omitempty" bson:"isFinished,omitempty"`
}

//EventQuery var is to query data from events collection
var EventQuery = database.MongoQuery{CollectionName: "events"}

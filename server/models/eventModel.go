package models

import (
	"go-react/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Event model
type Event struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name                  string             `json:"name" bson:"name"`
	Language              []string           `json:"language" bson:"language"`
	TypeID                primitive.ObjectID `json:"typeID" bson:"typeID"`
	Location              string             `json:"location" bson:"location"`
	Accommodation         string             `json:"accommodation" bson:"accommodation"`
	RegistrationCloseDate time.Time          `json:"registrationCloseDate" bson:"registrationCloseDate"`
	StartDate             time.Time          `json:"startDate" bson:"startDate"`
	EndDate               time.Time          `json:"endDate" bson:"endDate"`
	MaxParticipants       string             `json:"maxParticipants" bson:"maxParticipants"`
	Tags                  []string           `json:"tags" bson:"tags"`
	Description           string             `json:"description" bson:"description"`
	OwnerID               primitive.ObjectID `json:"ownerID" bson:"ownerID"`
	Budget                string             `json:"budget" bson:"budget"`
	IsApproved            bool               `json:"isApproved" bson:"isApproved"`
	Image                 string             `json:"image" bson:"image"`
	ReviewerID            primitive.ObjectID `json:"reviewerID" bson:"reviewerID"`
	IsFinished            bool               `json:"isFinished" bson:"isFinished"`
}

//EventQuery var is to query data from events collection
var EventQuery = database.MongoQuery{CollectionName: "events"}

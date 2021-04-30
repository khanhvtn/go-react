package models

import (
	"go-react/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User model
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

//UserQuery var is to query data from users collection
var UserQuery = database.MongoQuery{CollectionName: "users"}

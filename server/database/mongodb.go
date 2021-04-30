package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongoQuery struct is a collection fields for MongoDB
type MongoQuery struct {
	CollectionName string
}

//MongoClient var
var MongoClient *mongo.Client

//ConnectMongoDB func to connect to mongodb
func ConnectMongoDB() *mongo.Client {
	mongoURL := "mongodb+srv://khanhvtn93:khanhvtn93123@cluster0.zjom9.mongodb.net/goDB?authSource=admin&replicaSet=atlas-l3xb7s-shard-0&w=majority&readPreference=primary&appname=MongoDB%20Compass&retryWrites=true&ssl=true"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to MongoDB successfull.")
	MongoClient = client
	return client
}

//createCtxAndUserCol func is to create user collection, context, and cancel.
func createCtxAndUserCol(collectionName string) (col *mongo.Collection, ctx context.Context, cancel context.CancelFunc) {
	//get user collection
	col = MongoClient.Database("goDB").Collection(collectionName)
	//crete context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

//GetAll func is to return all record from a collection.
func (mongoQuery MongoQuery) GetAll() (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(mongoQuery.CollectionName)
	defer cancel()

	//create an empty array to store all fields from collection
	var data []bson.M

	//get all user record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	//map data to user variable
	if err = cur.All(ctx, &data); err != nil {
		return nil, err
	}
	//response data to client
	return data, nil
}

//GetOne func is to get one record from a collection
func (mongoQuery MongoQuery) GetOne(filter bson.M) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(mongoQuery.CollectionName)
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	result := bson.M{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err != mongo.ErrNoDocuments {
			//return err if there is a system error
			return nil, err
		}
		//return nil data when id is not existed.
		return nil, nil

	}

	return result, nil
}

//Create func is to create a new record to a collection
func (mongoQuery MongoQuery) Create(newData bson.M) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(mongoQuery.CollectionName)
	defer cancel()

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}
	newData["_id"] = insertResult.InsertedID
	return newData, nil
}

//UpdateOne func is to update one record from a collection
func (mongoQuery MongoQuery) UpdateOne(filter bson.M, update bson.M) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(mongoQuery.CollectionName)
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, nil
	}

	//query the new update
	newEventType, err := mongoQuery.GetOne(filter)
	if err != nil {
		return nil, err
	}

	return newEventType, nil
}

//DeleteOne func is to update one record from a collection
func (mongoQuery MongoQuery) DeleteOne(filter bson.M) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(mongoQuery.CollectionName)
	defer cancel()

	result, err := mongoQuery.GetOne(filter)
	if err != nil {
		return nil, err
	}

	//delete user from database
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		//response to client if there is an error.
		return nil, err
	}

	if deleteResult.DeletedCount == 0 {
		return nil, fiber.NewError(400, "Delete Fail.")
	}

	return result, nil
}

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

//MongGoQuery interface is a collections of method to query data in MongoDB
type MongoQuery interface {
	GetAll(collectionName string) (interface{}, error)
	GetOne(collectionName string) (interface{}, error)
	UpdateOne(collectionName string) (interface{}, error)
	DeleteOne(collectionName string) (interface{}, error)
}

//MongoDB struct is a collection fields for MongoDB
type MongoDB struct {
	Client *mongo.Client
}

//MongoClient var
var MongoClient MongoDB

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
	MongoClient.Client = client
	return client
}

//createCtxAndUserCol func is to create user collection, context, and cancel.
func createCtxAndUserCol(collectionName string) (col *mongo.Collection, ctx context.Context, cancel context.CancelFunc) {
	//get user collection
	col = MongoClient.Client.Database("goDB").Collection(collectionName)
	//crete context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

//GetAll func is to return all record from a collection.
func (mongoDB MongoDB) GetAll(collectionName string) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(collectionName)
	defer cancel()

	//create an empty array to store all fields from collection
	var data []bson.M

	//get all user record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fiber.NewError(500, "Something went wrong.")
	}
	defer cur.Close(ctx)
	//map data to user variable
	if err = cur.All(ctx, &data); err != nil {
		return nil, fiber.NewError(500, "Something went wrong.")
	}
	//response data to client
	return data, nil
}

//GetOne func is to get one record from a collection
func (mongoDB MongoDB) GetOne(collectionName string, id string) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(collectionName)
	defer cancel()

	//get id from client request
	idFilter, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		//response to client if there is an error.
		return nil, fiber.NewError(500, "Something went wrong.")
	}

	result := bson.M{}
	//Decode record into result
	if err := collection.FindOne(ctx, bson.M{"_id": idFilter}).Decode(&result); err != nil {
		if err != mongo.ErrNoDocuments {
			//return err if there is a system error
			return nil, fiber.NewError(500, "Something went wrong.")
		} else {
			//return nil data when id is not existed.
			return nil, nil
		}
	}
	return result, nil
}

//UpdateOne func is to update one record from a collection
func (mongoDB MongoDB) UpdateOne(collectionName string, filter bson.M, update bson.M) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(collectionName)
	defer cancel()

	//conver id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		id, err := primitive.ObjectIDFromHex(checkID.(string))
		if err != nil {
			return nil, fiber.NewError(500, "Something went wrong.")
		}
		filter["_id"] = id
	}

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		log.Fatal(err)
		return nil, fiber.NewError(500, "Something went wrong.")
	}
	if updateResult.MatchedCount == 0 {
		return nil, fiber.NewError(500, "Update Fail.")
	}
	return update, nil
}

//Delete func is to update one record from a collection
func (mongoDB MongoDB) DeleteOne(collectionName string, id string) (interface{}, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createCtxAndUserCol(collectionName)
	defer cancel()

	result, err := mongoDB.GetOne(collectionName, id)
	if err != nil {
		return nil, err
	}

	idFilter, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fiber.NewError(500, "Something went wrong.")
	}

	//delete user from database
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": idFilter})
	if err != nil {
		//response to client if there is an error.
		return nil, fiber.NewError(500, "Something went wrong.")
	}

	if deleteResult.DeletedCount == 0 {
		return nil, fiber.NewError(400, "Delete Fail.")
	}

	return result, nil
}

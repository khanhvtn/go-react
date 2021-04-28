package database

import (
	"context"
	"fmt"
	"go-react/utils"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongGoQuery interface is a collections of method to query data in MongoDB
type MongoQuery interface {
	GetAll(collectionName string) (interface{}, error)
	GetOne(collectionName string) (interface{}, error)
	UpdateOne(collectionName string) (interface{}, error)
	Delete(collectionName string) (interface{}, error)
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
	col, ctx, cancel := createCtxAndUserCol(collectionName)
	defer cancel()

	//create an empty array to store all fields from collection
	var users []bson.M

	//get all user record
	cur, err := userCol.Find(ctx, bson.D{})
	if err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}
	defer cur.Close(ctx)
	//map data to user variable
	if err = cur.All(ctx, &users); err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}
	//response data to client
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    users,
		Error:   nil})
}

package controllers

import (
	"context"
	"go-react/database"
	models "go-react/models"
	"go-react/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//createCtxAndUserCol func is to create user collection, context, and cancel.
func createCtxAndUserCol() (userCol *mongo.Collection, ctx context.Context, cancel context.CancelFunc) {
	//get user collection

	userCol = database.MongoClient.Client.Database("goDB").Collection("users")
	//crete context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

//GetUsers func to get all users.
func GetUsers(c *fiber.Ctx) error {
	//get user collection , context, cancel func
	userCol, ctx, cancel := createCtxAndUserCol()
	defer cancel()

	//create an empty array from user mdoel.
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

//CreateUser func to create a user.
func CreateUser(c *fiber.Ctx) error {
	//get user collection , context, cancel func
	userCol, ctx, cancel := createCtxAndUserCol()
	defer cancel()

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	//check user exists or not
	existedUser := new(models.User)
	if err := userCol.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existedUser); err != nil {
		if err != mongo.ErrNoDocuments {
			return fiber.NewError(500, "Something went wrong.")
		}
	} else {
		return c.Status(400).JSON(fiber.Map{"message": "Email already exists."})
	}

	//hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}
	//set hash password to new user
	user.Password = string(hashPassword)

	//create user in database
	insertResult, err := userCol.InsertOne(ctx, user)
	if err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data: fiber.Map{
			"_id":      insertResult.InsertedID,
			"email":    user.Email,
			"password": user.Password,
			"role":     user.Role,
		},
		Error: nil})
}

type userLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var jwtKey = []byte("jwtkey")

//Claims struct
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func genToken(user *models.User) (string, error) {
	expiredTime := time.Now().Add(60 * time.Second)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

//Login func to check user login
func Login(c *fiber.Ctx) error {
	//get user collection , context, cancel func
	userCol, ctx, cancel := createCtxAndUserCol()
	defer cancel()

	userLogin := new(userLogin)
	if err := c.BodyParser(userLogin); err != nil {
		return err
	}

	//check user exists or not
	existedUser := new(models.User)
	if err := userCol.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&existedUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(fiber.Map{"message": "Email or Password is invalid."})
		}
		return fiber.NewError(500, "Something went wrong.")
	}

	//check password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(userLogin.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email or Password is invalid."})
	}

	//create token
	token, err := genToken(existedUser)
	if err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}

	//send back token to client
	// return c.Status(200).JSON(fiber.Map{"token": token})
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    fiber.Map{"token": token},
		Error:   nil})
}

// UpdateUser func is to update user information
func UpdateUser(c *fiber.Ctx) error {
	//get user collection , context, cancel func
	userCol, ctx, cancel := createCtxAndUserCol()
	defer cancel()

	//get data client request
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Fatal(err)
		return fiber.NewError(500, "Something went wrong.")
	}

	// hash password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return fiber.NewError(500, "Something went wrong.")
	}
	user.Password = string(hashPass)

	//update user information
	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": user}
	updateResult, err := userCol.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return fiber.NewError(500, "Something went wrong.")
	}

	//response back to client
	// return c.Status(200).JSON(updateResult)
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    updateResult,
		Error:   nil})
}

//DeleteUser func is to delete an user.
func DeleteUser(c *fiber.Ctx) error {
	userCol, ctx, cancel := createCtxAndUserCol()
	defer cancel()

	//get id from client request
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		//response to client if there is an error.
		return fiber.NewError(500, "Something went wrong.")
	}
	//delete user from database
	deleteResult, err := userCol.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		//response to client if there is an error.
		return fiber.NewError(500, "Something went wrong.")
	}
	if deleteResult.DeletedCount == 0 {
		return fiber.NewError(400, "Invalid ID.")
	}
	//response to client when delete successful.
	// return c.Status(200).JSON(deleteResult)
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    deleteResult,
		Error:   nil})
}

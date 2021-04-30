package controllers

import (
	"context"
	"errors"
	"go-react/database"
	models "go-react/models"
	"go-react/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//createCtxAndUserCol func is to create user collection, context, and cancel.
func createCtxAndUserCol() (userCol *mongo.Collection, ctx context.Context, cancel context.CancelFunc) {
	//get user collection
	userCol = database.MongoClient.Database("goDB").Collection("users")
	//crete context with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

//GetUsers func to get all users.
func GetUsers(c *fiber.Ctx) error {
	//get all user record
	users, err := models.UserQuery.GetAll()
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})

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

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return err
	}

	//check user exists or not
	existedUser, err := models.UserQuery.GetOne(bson.M{"email": user.Email})
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	if existedUser != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email already exists."})
	}

	//hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	//set hash password to new user
	user.Password = string(hashPassword)

	//convert User type to bson.M
	bsonMUser, err := utils.InterfaceToBsonM(user)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}

	//create user in database
	newUser, err := models.UserQuery.Create(bsonMUser)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    newUser,
		Error:   nil})
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

func genToken(user *bson.M) (string, error) {
	expiredTime := time.Now().Add(60 * time.Second)
	claims := &Claims{
		Email: "email",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

//Login func to check user login
func Login(c *fiber.Ctx) error {

	userLogin := new(userLogin)
	if err := c.BodyParser(userLogin); err != nil {
		return err
	}

	//check user exists or not
	existedUser, err := models.UserQuery.GetOne(bson.M{"email": userLogin.Email})

	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	if existedUser == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email or Password is invalid."})
	}

	user := existedUser.(bson.M)

	//check password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(userLogin.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email or Password is invalid."})
	}

	//create token
	token, err := genToken(&user)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}

	//send back token to client
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    fiber.Map{"token": token},
		Error:   nil})
}

// UpdateUser func is to update user information
func UpdateUser(c *fiber.Ctx) error {
	//get id param
	id := c.Params("id")
	//get data client request
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}

	// hash password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	user.Password = string(hashPass)

	//convert User type to bson.M
	bsonMUser, err := utils.InterfaceToBsonM(user)

	//update user information
	filter := bson.M{"_id": id}

	updateResult, err := models.UserQuery.UpdateOne(filter, bsonMUser)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}

	if updateResult == nil {
		return fiber.NewError(400, "Update Fail.")
	}

	//response back to client
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    updateResult,
		Error:   nil})
}

//DeleteUser func is to delete an user.
func DeleteUser(c *fiber.Ctx) error {

	//get id from client request
	id := c.Params("id")
	//delete user from database
	deletedResult, err := models.UserQuery.DeleteOne(bson.M{"_id": id})
	if err != nil {
		//response to client if there is an error.
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errors.New("Something went wrong")})
	}
	if deletedResult == nil {
		return fiber.NewError(400, "Delete Fail.")
	}
	//response to client when delete successful.
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    deletedResult,
		Error:   nil})
}

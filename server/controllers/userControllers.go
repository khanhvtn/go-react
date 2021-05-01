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
	users, errQuery := models.UserQuery.GetAll()
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
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
	existedUser, errQuery := models.UserQuery.GetOne(bson.M{"email": user.Email})
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
	}
	if existedUser != nil {
		return c.Status(400).JSON(fiber.Map{"message": ""})
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    400,
			Data:    nil,
			Error:   errors.New("Email already exists")})
	}

	//hash password
	hashPassword, errBcrypt := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if errBcrypt != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errBcrypt})
	}
	//set hash password to new user
	user.Password = string(hashPassword)

	//convert User type to bson.M
	bsonMUser, errConvert := utils.InterfaceToBsonM(user)
	if errConvert != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errConvert})
	}

	//create user in database
	newUser, err := models.UserQuery.Create(bsonMUser)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    err.Code,
			Data:    nil,
			Error:   err})
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
	existedUser, errQuery := models.UserQuery.GetOne(bson.M{"email": userLogin.Email})

	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
	}
	if len(existedUser.(bson.M)) == 0 {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    400,
			Data:    nil,
			Error:   errors.New("Email or Password is invalid")})

	}

	user := existedUser.(bson.M)

	//check password is valid
	if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(userLogin.Password)); err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    400,
			Data:    nil,
			Error:   errors.New("Email or Password is invalid")})
	}

	//create token
	token, errGenToken := genToken(&user)
	if errGenToken != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   errGenToken})
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
			Error:   err})
	}

	// hash password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}
	user.Password = string(hashPass)

	//convert User type to bson.M
	bsonMUser, err := utils.InterfaceToBsonM(user)

	//update user information
	filter := bson.M{"_id": id}

	updateResult, errQuery := models.UserQuery.UpdateOne(filter, bsonMUser)
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
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
	deletedResult, errQuery := models.UserQuery.DeleteOne(bson.M{"_id": id})
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
	}
	//response to client when delete successful.
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    deletedResult,
		Error:   nil})
}

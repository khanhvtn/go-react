package controllers

import (
	"errors"
	"go-react/models"
	"go-react/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func genDefaultError(c *fiber.Ctx) error {
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    500,
		Data:    nil,
		Error:   errors.New("Something went wrong")})
}

//GetEventTypes func is to return all event types.
func GetEventTypes(c *fiber.Ctx) error {
	eventTypes, err := models.EventTypeQuery.GetAll()
	if err != nil {
		return genDefaultError(c)
	}

	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    eventTypes,
		Error:   nil})

}

//GetEventType func is to return an event type.
func GetEventType(c *fiber.Ctx) error {
	id := c.Params("id")

	eventType, err := models.EventTypeQuery.GetOne(bson.M{"_id": id})
	if err != nil {
		return genDefaultError(c)

	}
	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    eventType,
		Error:   nil})
}

//CreateEventType func is to create an event type.
func CreateEventType(c *fiber.Ctx) error {
	userReq := new(models.EventType)

	if err := c.BodyParser(userReq); err != nil {
		return genDefaultError(c)
	}

	existedEventTypeChan := make(chan bson.M)

	go func() {
		//check event type exists or not
		existedEventType, err := models.EventTypeQuery.GetOne(bson.M{"name": userReq.Name})
		existedEventTypeChan <- bson.M{
			"eventType": existedEventType,
			"error":     err,
		}
	}()

	//convert  userReq to bson.M
	bsonMEventType, err := utils.InterfaceToBsonM(userReq)
	if err != nil {
		return genDefaultError(c)
	}
	//get result from existedEventTypeChan
	result := <-existedEventTypeChan
	if result["error"] != nil {
		return genDefaultError(c)
	}

	if result["eventType"] != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    400,
			Data:    nil,
			Error:   errors.New("Event Type already existed")})
	}

	//Create a new event typoe
	newEventType, err := models.EventTypeQuery.Create(bsonMEventType)
	if err != nil {
		return genDefaultError(c)
	}

	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    newEventType,
		Error:   nil})

}

//UpdateEventType func is to update an event type.
func UpdateEventType(c *fiber.Ctx) error {
	id := c.Params("id")
	userReq := new(models.EventType)

	if err := c.BodyParser(userReq); err != nil {
		return genDefaultError(c)
	}

	//convert  userReq to bson.M
	bsonMEventType, err := utils.InterfaceToBsonM(userReq)
	if err != nil {
		return genDefaultError(c)
	}

	//filter
	filter := bson.M{"_id": id}

	//Update Event Type
	updatedEventType, err := models.EventTypeQuery.UpdateOne(filter, bsonMEventType)
	if err != nil {
		return genDefaultError(c)
	}

	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    updatedEventType,
		Error:   nil})
}

//DeleteEventType func is to delete an event type.
func DeleteEventType(c *fiber.Ctx) error {
	id := c.Params("id")

	//delete an event type
	deletedEventType, err := models.EventTypeQuery.DeleteOne(bson.M{"_id": id})

	if err != nil {
		return genDefaultError(c)
	}

	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    deletedEventType,
		Error:   nil})

}

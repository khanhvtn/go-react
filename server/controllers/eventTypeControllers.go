package controllers

import (
	"errors"
	"go-react/models"
	"go-react/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

//GetEventTypes func is to return all event types.
func GetEventTypes(c *fiber.Ctx) error {
	eventTypes, errQuery := models.EventTypeQuery.GetAll()
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
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

	eventType, errQuery := models.EventTypeQuery.GetOne(bson.M{"_id": id})
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
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
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}

	existedEventTypeChan := make(chan bson.M)

	go func() {
		//check event type exists or not
		existedEventType, err := models.EventTypeQuery.GetOne(bson.M{"name": userReq.Name})
		if err != nil {
			existedEventTypeChan <- bson.M{
				"eventType": existedEventType,
				"error": bson.M{
					"message": err.Message,
					"code":    err.Code,
				},
			}
		} else {
			existedEventTypeChan <- bson.M{
				"eventType": existedEventType,
				"error":     nil,
			}
		}
	}()

	//convert  userReq to bson.M
	bsonMEventType, err := utils.InterfaceToBsonM(userReq)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}
	//get result from existedEventTypeChan
	result := <-existedEventTypeChan
	if result["error"] != nil {
		err := result["error"].(bson.M)
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    err["code"].(int),
			Data:    nil,
			Error:   errors.New(err["message"].(string))})
	}

	if len(result["eventType"].(bson.M)) != 0 {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    400,
			Data:    nil,
			Error:   errors.New("Event Type already existed")})
	}

	//Create a new event typoe
	newEventType, errQuery := models.EventTypeQuery.Create(bsonMEventType)
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errQuery})
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
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}

	//convert  userReq to bson.M
	bsonMEventType, err := utils.InterfaceToBsonM(userReq)
	if err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}

	//filter
	filter := bson.M{"_id": id}

	//Update Event Type
	updatedEventType, errQuery := models.EventTypeQuery.UpdateOne(filter, bsonMEventType)
	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
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
	deletedEventType, errQuery := models.EventTypeQuery.DeleteOne(bson.M{"_id": id})

	if errQuery != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    errQuery.Code,
			Data:    nil,
			Error:   errors.New(errQuery.Message)})
	}

	return utils.CusResponse(utils.CusResp{
		Context: c,
		Code:    200,
		Data:    deletedEventType,
		Error:   nil})

}

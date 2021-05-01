package controllers

import (
	"errors"
	"go-react/models"
	"go-react/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

//GetEvents func is to return all Events.
func GetEvents(c *fiber.Ctx) error {
	events, errQuery := models.EventQuery.GetAll()
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
		Data:    events,
		Error:   nil})

}

//GetEvent func is to return an Event.
func GetEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	event, errQuery := models.EventQuery.GetOne(bson.M{"_id": id})
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
		Data:    event,
		Error:   nil})

}

//CreateEvent func is to create an Event.
func CreateEvent(c *fiber.Ctx) error {
	var userReq bson.M

	if err := c.BodyParser(userReq); err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}

	newEvent, errQuery := models.EventQuery.Create(userReq)
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
		Data:    newEvent,
		Error:   nil})
}

//UpdateEvent func is to update an Event.
func UpdateEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	var userReq bson.M
	if err := c.BodyParser(userReq); err != nil {
		return utils.CusResponse(utils.CusResp{
			Context: c,
			Code:    500,
			Data:    nil,
			Error:   err})
	}
	filter := bson.M{"_id": id}
	updatedEvent, errQuery := models.EventQuery.UpdateOne(filter, userReq)

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
		Data:    updatedEvent,
		Error:   nil})

}

//DeleteEvent func is to delete an Event.
func DeleteEvent(c *fiber.Ctx) error {
	id := c.Params("id")
	filter := bson.M{"_id": id}
	deletedEvent, errQuery := models.EventQuery.DeleteOne(filter)
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
		Data:    deletedEvent,
		Error:   nil})

}

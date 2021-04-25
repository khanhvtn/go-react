package utils

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

//CusResp struct
type CusResp struct {
	Context *fiber.Ctx
	Code    int
	Data    interface{}
	Error   error
}

//CusResponse func is to customize response and send it to client.
func CusResponse(cusResp CusResp) error {
	if cusResp.Error != nil {
		return cusResp.Context.Status(cusResp.Code).JSON(fiber.Map{
			"code":       cusResp.Code,
			"message":    "failure",
			"errMessage": cusResp.Error.Error(),
		})
	}
	//calcualte len data
	var dataLen int

	if result := reflect.TypeOf(cusResp.Data); result.Kind() == reflect.Slice {
		fmt.Println(result.Kind().String())
		castData := cusResp.Data.([]bson.M)
		dataLen = len(castData)
	} else {
		dataLen = 1
	}

	return cusResp.Context.Status(cusResp.Code).JSON(fiber.Map{
		"code":       cusResp.Code,
		"message":    "success",
		"dataLength": dataLen,
		"data":       cusResp.Data,
	})
}

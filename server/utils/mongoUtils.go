package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

//InterfaceToBsonM func is to convert an interface to bson.M
func InterfaceToBsonM(val interface{}) (bson.M, error) {
	var result bson.M
	byteInterface, err := bson.Marshal(val)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(byteInterface, &result); err != nil {
		return nil, err
	}
	return result, nil
}

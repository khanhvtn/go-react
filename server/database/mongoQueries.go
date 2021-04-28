package database

//MongGoQuery interface is a collections of method to query data in MongoDB
type MongGoQuery interface {
	GetAll() (interface{}, error)
	GetOne() (interface{}, error)
	UpdateOne() (interface{}, error)
	Delete() (interface{}, error)
}

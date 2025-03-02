package dtos

type Database struct {
	Id     int64
	Name   string
	Server string
}

type NewDatabase struct {
	Name   string
	Server string
}

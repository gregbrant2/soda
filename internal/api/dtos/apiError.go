package dtos

type ApiError struct {
	Message string
	Fields  map[string]string
}

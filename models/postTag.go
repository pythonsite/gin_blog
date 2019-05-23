package models

type PostTag struct {
	BaseModel
	PostId uint
	TagId uint
}
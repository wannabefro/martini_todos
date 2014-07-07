package models

type Todo struct {
	Id					int64
	Created			int64
	Title				string `form:"Title" binding:"required"`
	Description	string `form:"Description" binding:"required`
}

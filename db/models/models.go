package models

type Todo struct {
	Id					int64
	Created			int64
	Title				string `form:"Title"`
	Description	string `form:"Description"`
}

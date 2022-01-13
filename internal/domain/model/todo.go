package model

type Todo struct {
	UserId    int64  `json:",omitempty"`
	Id        int64  `json:",omitempty"`
	Title     string `json:",omitempty"`
	Completed bool   `json:",omitempty"`
}

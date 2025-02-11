package model

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Balance  int    `json:"balance"`
	Password string `json:"password"`
}

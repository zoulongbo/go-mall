package models

type User struct {
	ID           int64  `json:"id" sql:"id" info:"id"`
	Nickname     string `json:"nickname" sql:"nickname" info:"nickname"`
	Username     string `json:"username" sql:"username" info:"username"`
	HashPassword string `json:"-" sql:"password" info:"password"`
}

const (
	UserTable = "user"
)

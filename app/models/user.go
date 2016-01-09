package models

import (
	"github.com/pborman/uuid"
)

type User struct {
	Id               string
	KoofrAccessToken string
}

var db = map[string]*User{}

func GetUser(id string) *User {
	return db[id]
}

func NewUser() *User {
	user := &User{
		Id: uuid.New(),
	}
	db[user.Id] = user
	return user
}

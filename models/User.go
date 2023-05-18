package models

import (
	"otus-sonet/utils"
	"time"
)

// User structure in DB
type User struct {
	Id         string `json:"id"          sql:"id"`
	FirstName  string `json:"first_name"  sql:"first_name"`
	SecondName string `json:"second_name" sql:"second_name"`
	Age        int    `json:"age"         sql:"-"`
	Birthdate  string `json:"birthdate"   sql:"birthdate"`
	Biography  string `json:"biography"   sql:"biography"`
	City       string `json:"city"        sql:"city"`
	Password   string `json:"-"           sql:"password"`
	Token      string `json:"-"           sql:"token"`
	TokenTill  string `json:"-"           sql:"token_till"`
}

func CalcUserAge(u *User) {
	birthday, err := time.Parse("2006-01-02", u.Birthdate)
	if err == nil {
		age, _, _, _, _, _ := utils.TimeDiff(birthday, time.Now())
		u.Age = age
	}
}

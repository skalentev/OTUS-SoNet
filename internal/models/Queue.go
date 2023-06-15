package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Queue struct {
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

func QueueSubscribe(ctx context.Context, rc *redis.Client, qKey string) {
	for {
		qLen := rc.LLen(ctx, qKey).Val()
		if qLen > 0 {
			msg := rc.LPop(ctx, qKey).Val()
			fmt.Println("Received message: " + msg)
		} else {
			time.Sleep(300 * time.Millisecond)
		}
	}
}

package models

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Friend struct {
	redis  *redis.Client
	db     *sql.DB
	router *gin.Engine
}

var Friends = &Friend{}

func (f *Friend) Init(redis *redis.Client, db *sql.DB) {
	f.redis = redis
	f.db = db
}

func (f *Friend) Set(userId string, friendId string) error {
	id := uuid.NewString()
	var query string
	switch DB.Driver {
	case "mysql":
		query = "INSERT INTO friend SET id = ?, `user_id` = ?, `friend_id` = ?"
	default:
		query = "INSERT INTO public.friend ( id, user_id, friend_id) VALUES ($1, $2, $3)"
	}
	_, err := DB.DB.Exec(query, id, userId, friendId)
	return err
}

func (f *Friend) Delete(userId string, friendId string) error {
	var query string
	switch DB.Driver {
	case "mysql":
		query = "DELETE FROM friend WHERE `user_id` = ? AND `friend_id` = ?"
	default:
		query = "DELETE FROM public.friend WHERE user_id = $1 AND friend_id = $2"
	}
	_, err := DB.DB.Exec(query, userId, friendId)
	return err
}

func (f *Friend) Get(userId string) ([]string, error) {
	return nil, nil
}

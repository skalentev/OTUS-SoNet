package models

import (
	"database/sql"
	"fmt"
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
	var query string
	switch DBSlave.Driver {
	case "mysql":
		query = "SELECT f.friend_id from friend f WHERE f.user_id = ? ORDER BY f.friend_id "
	default:
		query = "SELECT f.friend_id from public.friend f WHERE f.user_id = $1 ORDER BY f.friend_id"
	}
	rows, err := DBSlave.DB.Query(query, userId)
	if err != nil {
		fmt.Println("DBErr:", err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("Row close error:", err)
			return
		}
	}()

	var friends []string

	for rows.Next() {
		var friend string
		if err := rows.Scan(&friend); err != nil {
			fmt.Println("Next err:", err)
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

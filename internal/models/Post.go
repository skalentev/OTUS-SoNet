package models

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

/*
	{
	  "id": "1d535fd6-7521-4cb1-aa6d-031be7123c4d",
	  "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lectus mauris ultrices eros in cursus turpis massa.",
	  "author_user_id": "string"
	}
*/
type Post struct {
	Id     string `json:"id"              sql:"id"`
	Text   string `json:"text"            sql:"text"`
	UserId string `json:"author_user_id"  sql:"user_id"`
}

type TPost struct {
	redis  *redis.Client
	db     *sql.DB
	router *gin.Engine
}

var Posts = &TPost{}

func (p *TPost) Init(redis *redis.Client, db *sql.DB) {
	p.redis = redis
	p.db = db
}

// Create Добавляет пост в хранилище
func (p *TPost) Create(post Post) (string, error) {
	id := uuid.NewString()
	var query string
	switch DB.Driver {
	case "mysql":
		query = "INSERT INTO post SET id = ?, `text` = ?, `user_id` = ? "
	default:
		query = "INSERT INTO public.post ( id, text, user_id) VALUES ($1, $2, $3)"
	}
	_, err := DB.DB.Exec(query, id, post.Text, post.UserId)
	if err != nil {
		fmt.Println(query)
		return "", err
	}
	return id, nil
}

// Update Изменение поста по id и пользователю
func (p *TPost) Update(post Post) error {
	var query string
	switch DB.Driver {
	case "mysql":
		query = "UPDATE post SET `text` = ? WHERE id = ? AND `user_id` = ?"
	default:
		query = "UPDATE public.post SET text= $1 WHERE id = $2 AND user_id = $3"
	}
	_, err := DB.DB.Exec(query, post.Text, post.Id, post.UserId)
	if err != nil {
		fmt.Println(query)
		return err
	}
	return nil
}

// Delete Удаление поста по id пользователя и id поста
func (p *TPost) Delete(postId string, userId string) error {
	var query string
	switch DB.Driver {
	case "mysql":
		query = "DELETE FROM post WHERE `id` = ? AND `user_id` = ?"
	default:
		query = "DELETE FROM public.post WHERE id = $1 AND user_id = $2"
	}
	_, err := DB.DB.Exec(query, postId, userId)
	return err
}

func (p *TPost) Get(postId string) (*Post, error) {
	var post Post
	var query string
	switch DBSlave.Driver {
	case "mysql":
		query = "SELECT p.id, p.text, p.user_id from post p WHERE p.id = ? LIMIT 1"
	default:
		query = "SELECT id, text, user_id from public.post WHERE id = $1 limit 1"
	}
	if err := DBSlave.DB.QueryRow(query, postId).Scan(&post.Id, &post.Text, &post.UserId); err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	fmt.Println("post:", post)

	if _, err := uuid.Parse(post.Id); err != nil {
		return nil, err
	}
	return &post, nil

}

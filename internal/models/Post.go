package models

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

func (p *TPost) Create(j *Post) (string, error) {
	return "", nil
}

func (p *TPost) Update(j *Post) error {
	return nil
}
func (p *TPost) Delete(j string) error {
	return nil
}
func (p *TPost) Get(j string) (*Post, error) {
	return nil, nil
}

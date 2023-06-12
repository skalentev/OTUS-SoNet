package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	models2 "otus-sonet/internal/models"
	utils2 "otus-sonet/internal/utils"
	"time"
)

func UserRegister(c *gin.Context) {
	var user models2.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(400)
		return
	}

	if birthday, err := time.Parse("2006-01-02", user.Birthdate); err != nil {
		c.AbortWithStatus(400)
		return
	} else {
		user.Birthdate = birthday.Format("2006-01-02")
	}

	var errHash error
	user.Password, errHash = utils2.GenerateHashPassword(user.Password)
	if errHash != nil {
		utils2.Code500(c, "could not generate password hash", -4)
		return
	}
	user.Id = utils2.GenerateToken()

	var query string
	switch models2.DB.Driver {
	case "mysql":
		query = "INSERT INTO user SET id = ?, `first_name` = ?, `second_name` = ?, `birthdate` = ?, `city` = ?, `biography` = ?, `password` = ? "
	default:
		query = "INSERT INTO public.user ( id, first_name, second_name, birthdate, city, biography, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	}
	_, err := models2.DB.DB.Exec(query,
		user.Id, user.FirstName, user.SecondName, user.Birthdate, user.City, user.Biography, user.Password)
	if err != nil {
		fmt.Println(query)
		utils2.Code500(c, err.Error(), -5)
		return
	}
	c.JSON(200, gin.H{"user_id": user.Id})
}

func UserGetId(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(400)
		return
	}

	start := time.Now()
	var user models2.User
	var query string
	switch models2.DBRO.Driver {
	case "mysql":
		query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from user u WHERE u.id = ? LIMIT 1"
	default:
		query = "SELECT id, first_name, second_name, birthdate, biography, city from public.user WHERE id = $1 limit 1"
	}
	if err := models2.DBRO.DB.QueryRow(query,
		id).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(query)
			c.AbortWithStatus(404)
			return
		} else {
			utils2.Code500(c, "service unavailable", -6)
			fmt.Println("err:", err)
			return
		}
	}
	models2.Prom.DbTimeSummary.WithLabelValues("select", "userGetId", "query").Observe(float64(time.Since(start).Milliseconds()))
	models2.Prom.DbTimeGauge.WithLabelValues("select", "userGetId", "query").Set(float64(time.Since(start).Milliseconds()))
	models2.CalcUserAge(&user)

	fmt.Println("user:", user)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, user)
}

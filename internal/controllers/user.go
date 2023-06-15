package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"otus-sonet/internal/models"
	"otus-sonet/internal/utils"
	"time"
)

func UserRegister(c *gin.Context) {
	var user models.User

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
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		utils.Code500(c, "could not generate password hash", -4)
		return
	}
	user.Id = utils.GenerateToken()

	var query string
	switch models.DB.Driver {
	case "mysql":
		query = "INSERT INTO user SET id = ?, `first_name` = ?, `second_name` = ?, `birthdate` = ?, `city` = ?, `biography` = ?, `password` = ? "
	default:
		query = "INSERT INTO public.user ( id, first_name, second_name, birthdate, city, biography, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	}
	_, err := models.DB.DB.Exec(query,
		user.Id, user.FirstName, user.SecondName, user.Birthdate, user.City, user.Biography, user.Password)
	if err != nil {
		fmt.Println(query)
		utils.Code500(c, err.Error(), -5)
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
	var user models.User
	var query string
	switch models.DBSlave.Driver {
	case "mysql":
		query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from user u WHERE u.id = ? LIMIT 1"
	default:
		query = "SELECT id, first_name, second_name, birthdate, biography, city from public.user WHERE id = $1 limit 1"
	}
	if err := models.DBSlave.DB.QueryRow(query,
		id).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(query)
			c.AbortWithStatus(404)
			return
		} else {
			utils.Code500(c, "service unavailable", -6)
			fmt.Println("err:", err)
			return
		}
	}
	models.Prom.DbTimeSummary.WithLabelValues("select", "userGetId", "query").Observe(float64(time.Since(start).Milliseconds()))
	models.Prom.DbTimeGauge.WithLabelValues("select", "userGetId", "query").Set(float64(time.Since(start).Milliseconds()))
	models.CalcUserAge(&user)

	fmt.Println("user:", user)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, user)
}

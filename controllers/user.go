package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"otus-sonet/models"
	"otus-sonet/utils"
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

	var password string
	if err := models.DB.QueryRow("SELECT password from user WHERE id = ? LIMIT 1", user.Id).Scan(&password); err == nil {
		c.AbortWithStatus(400)
		fmt.Println("user already exists")
		return
	} else if err != sql.ErrNoRows {
		utils.Code500(c, "service unavailable", -3)
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		utils.Code500(c, "could not generate password hash", -4)
		return
	}
	user.Id = utils.GenerateToken()
	_, err := models.DB.Exec("INSERT INTO user SET id = ?, `first_name` = ?, `second_name` = ?, `birthdate` = ?, "+
		"`city` = ?, `biography` = ?, `password` = ? ",
		user.Id, user.FirstName, user.SecondName, user.Birthdate, user.City, user.Biography, user.Password)
	if err != nil {
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

	var user models.User
	if err := models.DB.QueryRow("SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from user u WHERE u.id = ? LIMIT 1",
		id).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(404)
			return
		} else {
			utils.Code500(c, "service unavailable", -6)
			fmt.Println("err:", err)
			return
		}
	}
	models.CalcUserAge(&user)

	fmt.Println("user:", user)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(200, user)
}

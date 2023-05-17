package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-auth/models"
	"go-auth/utils"
	"time"
)

func UserRegister(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if birthday, err := time.Parse(time.DateOnly, user.Birthdate); err != nil {
		c.JSON(400, gin.H{"error": "data error (" + user.Birthdate + ")"})
		return
	} else {
		user.Birthdate = birthday.Format(time.DateOnly)
	}

	var password string
	if err := models.DB.QueryRow("SELECT password from user WHERE id = ? LIMIT 1", user.Id).Scan(&password); err == nil {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(500, gin.H{"error": "service unavailable"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}
	user.Id = utils.GenerateToken()
	_, err := models.DB.Exec("INSERT INTO user SET id = ?, `first_name` = ?, `second_name` = ?, `birthdate` = ?, "+
		"`city` = ?, `biography` = ?, `password` = ? ",
		user.Id, user.FirstName, user.SecondName, user.Birthdate, user.City, user.Biography, user.Password)
	if err != nil {
		c.JSON(500, gin.H{"message": err, "code": -5})
		return
	}
	c.JSON(200, gin.H{"user_id": user.Id})
}

func UserGetId(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(400, "")
		return
	}

	var user models.User
	if err := models.DB.QueryRow("SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from user u WHERE u.id = ? LIMIT 1",
		id).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, "")
			return
		} else {
			c.JSON(500, gin.H{"message": "service unavailable", "request_id": utils.GetRequestId(c), "code": -1})
			fmt.Println("err:", err)
			return
		}
	}
	models.CalcUserAge(&user)

	fmt.Println("user:", user)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.JSON(404, "")
		return
	}

	c.JSON(200, user)
}

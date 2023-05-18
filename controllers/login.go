package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"otus-sonet/models"
	"otus-sonet/utils"
	"time"
)

func Login(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{})
		return
	}

	var password string

	if err := models.DB.QueryRow("SELECT password from user WHERE id = ? LIMIT 1", user.Id).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{})
			return
		} else {
			c.JSON(500, gin.H{"message": "service unavailable", "request_id": utils.GetRequestId(c), "code": -1})
			return
		}
	}

	errHash := utils.CompareHashPassword(user.Password, password)
	if !errHash {
		c.JSON(400, gin.H{})
		return
	}

	tokenTime := time.Now().Add(5 * time.Minute)
	token := utils.GenerateToken()

	_, err := models.DB.Exec("INSERT INTO session SET token = ?, `user_id` = ?, `token_till` = ? ",
		token, user.Id, tokenTime)
	if err != nil {
		c.JSON(500, gin.H{"message": "could not save session", "request_id": utils.GetRequestId(c), "code": -2})
		return
	}

	//if err != nil {
	//	c.JSON(500, gin.H{"error": "could not generate token"})
	//	return
	//}

	//	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"token": token})
}

func Home(c *gin.Context) {

	c.JSON(200, gin.H{"success": "home page", "requestId": utils.GetRequestId(c)})
}

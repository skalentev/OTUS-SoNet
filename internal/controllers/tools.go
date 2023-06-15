package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"otus-sonet/internal/models"
	"otus-sonet/internal/utils"
)

func GetAuthUser(c *gin.Context) (models.User, error) {

	val, exists := c.Get("user")
	if !exists {
		return models.User{}, errors.New("no auth key")
	}

	var user models.User = val.(models.User)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.AbortWithStatus(401)
		return models.User{}, errors.New("no auth uuid")
	}

	return user, nil
}

func AuthTest(c *gin.Context) {

	user, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "OK"})
}

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "no route or wrong method, try GET /user/get/{id}", "requestId": utils.GetRequestId(c), "code": 0})
}

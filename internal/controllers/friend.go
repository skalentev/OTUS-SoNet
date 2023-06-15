package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"otus-sonet/internal/models"
	"otus-sonet/internal/utils"
)

func FriendSet(c *gin.Context) {

	friendId := c.Param("user_id")
	if friendId == "" {
		c.AbortWithStatus(400)
		return
	}

	if _, err := uuid.Parse(friendId); err != nil {
		c.AbortWithStatus(400)
		return
	}

	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if err := models.Friends.Set(userId.Id, friendId); err != nil {
		utils.Code500(c, err.Error(), -5)
		return
	}
	c.AbortWithStatus(200)
}

func FriendDelete(c *gin.Context) {

	friendId := c.Param("user_id")
	if friendId == "" {
		c.AbortWithStatus(400)
		return
	}

	if _, err := uuid.Parse(friendId); err != nil {
		c.AbortWithStatus(400)
		return
	}

	userId, err := GetAuthUser(c)
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if err := models.Friends.Delete(userId.Id, friendId); err != nil {
		utils.Code500(c, err.Error(), -5)
		return
	}
	c.AbortWithStatus(200)
}

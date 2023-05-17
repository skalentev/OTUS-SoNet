package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-auth/models"
)

func Test(c *gin.Context) {

	val, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User = val.(models.User)

	fmt.Println("user:", user)

	if _, err := uuid.Parse(user.Id); err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(200, gin.H{"user": user})
}

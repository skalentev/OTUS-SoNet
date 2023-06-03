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
		c.AbortWithStatus(400)
		return
	}

	var password string
	var query string
	switch models.Driver {
	case "mysql":
		query = "SELECT password from user WHERE id = ? LIMIT 1"
	default:
		query = "SELECT password from public.user WHERE id = $1 limit 1"
	}
	if err := models.DB.QueryRow(query, user.Id).Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(404)
			return
		} else {
			utils.Code500(c, "Service unavailable", -1)
			return
		}
	}

	errHash := utils.CompareHashPassword(user.Password, password)
	if !errHash {
		c.AbortWithStatus(404)
		return
	}

	tokenTime := time.Now().Add(5 * time.Minute)
	token := utils.GenerateToken()
	switch models.Driver {
	case "mysql":
		query = "INSERT INTO session SET token = ?, `user_id` = ?, `token_till` = ? "
	default:
		query = "INSERT INTO public.session ( token, user_id, token_till) VALUES ($1, $2, $3)"
	}
	_, err := models.DB.Exec(query, token, user.Id, tokenTime)
	if err != nil {
		utils.Code500(c, "Could not save session: "+err.Error(), -2)
		return
	}

	c.JSON(200, gin.H{"token": token})
}

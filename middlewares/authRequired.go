package middlewares

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go-auth/models"
	"go-auth/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if headerParts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, err := uuid.Parse(headerParts[1]); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := models.DB.QueryRow("SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from session s, user u WHERE s.token = ? AND s.token_till>? AND u.id=s.user_id LIMIT 1",
			headerParts[1], time.Now()).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{})
				return
			} else {
				c.JSON(500, gin.H{"message": "service unavailable", "request_id": utils.GetRequestId(c), "code": -1})
				fmt.Println("err:", err)
				return
			}
		}
		models.CalcUserAge(&user)
		c.Set("user", user)
		c.Next()
	}
}

package middlewares

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"otus-sonet/models"
	"otus-sonet/utils"
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

		var query string
		switch models.Driver {
		case "mysql":
			query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from session s, user u WHERE s.token = ? AND s.token_till>? AND u.id=s.user_id LIMIT 1"
		default:
			query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, u.biography, u.city from public.session s, public.user u WHERE s.token = $1 AND s.token_till>$2 AND u.id=s.user_id limit 1"
		}
		if err := models.DB.QueryRow(query,
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

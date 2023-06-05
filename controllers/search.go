package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"otus-sonet/models"
	"otus-sonet/utils"
	"time"
)

func UserSearch(c *gin.Context) {

	firstName := c.Query("first_name")
	lastName := c.Query("last_name")
	if firstName == "" {
		c.AbortWithStatus(400)
		fmt.Println("No firstName")
		return
	}

	if lastName == "" {
		c.AbortWithStatus(400)
		fmt.Println("No lastName")
		return
	}

	start := time.Now()
	var query string
	switch models.DBRO.Driver {
	case "mysql":
		query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, COALESCE(u.biography,'-') as biography, u.city from user u WHERE u.first_name LIKE ? AND u.second_name LIKE ? ORDER BY u.id "
	default:
		query = "SELECT id, first_name, second_name, birthdate, COALESCE(biography,'-') as biography, city from public.user WHERE first_name LIKE $1 AND second_name LIKE $2 ORDER BY id"
	}
	rows, err := models.DBRO.DB.Query(query,
		firstName+"%", lastName+"%")
	if err != nil {
		utils.Code500(c, "Query error", -7)
		fmt.Println("DBErr:", err)
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println("Row close error:", err)
			return
		}
	}()
	models.Prom.DbTimeSummary.WithLabelValues("select", "userSearch", "query").Observe(float64(time.Since(start).Milliseconds()))
	models.Prom.DbTimeGauge.WithLabelValues("select", "userSearch", "query").Set(float64(time.Since(start).Milliseconds()))

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
			utils.Code500(c, "get data error", -8)
			fmt.Println(err)
			return
		}
		models.CalcUserAge(&user)
		users = append(users, user)
	}
	models.Prom.DbTimeSummary.WithLabelValues("select", "userSearch", "rows").Observe(float64(time.Since(start).Milliseconds()))
	if err = rows.Err(); err != nil {
		utils.Code500(c, "Metrics error", -9)
		return
	}

	c.JSON(200, users)
}

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models2 "otus-sonet/internal/models"
	"otus-sonet/internal/utils"
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
	switch models2.DBRO.Driver {
	case "mysql":
		query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, COALESCE(u.biography,'-') as biography, u.city from user u WHERE u.first_name LIKE ? AND u.second_name LIKE ? ORDER BY u.id "
	default:
		query = "SELECT u.id, u.first_name, u.second_name, u.birthdate, COALESCE(u.biography,'-') as biography, u.city from public.user u WHERE u.first_name LIKE $1 AND u.second_name LIKE $2 ORDER BY u.id"
	}
	rows, err := models2.DBRO.DB.Query(query,
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
	models2.Prom.DbTimeSummary.WithLabelValues("select", "userSearch", "query").Observe(float64(time.Since(start).Milliseconds()))
	models2.Prom.DbTimeGauge.WithLabelValues("select", "userSearch", "query").Set(float64(time.Since(start).Milliseconds()))

	var users []models2.User

	for rows.Next() {
		var user models2.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City); err != nil {
			utils.Code500(c, "get data error", -8)
			fmt.Println(err)
			return
		}
		models2.CalcUserAge(&user)
		users = append(users, user)
	}
	models2.Prom.DbTimeSummary.WithLabelValues("select", "userSearch", "rows").Observe(float64(time.Since(start).Milliseconds()))
	if err = rows.Err(); err != nil {
		utils.Code500(c, "Metrics error", -9)
		return
	}

	c.JSON(200, users)
}

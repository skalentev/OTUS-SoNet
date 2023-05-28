package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"otus-sonet/models"
	"otus-sonet/utils"
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

	rows, err := models.DB.Query("SELECT u.id, u.first_name, u.second_name, u.birthdate, COALESCE(u.biography,'-') as biography, u.city from user u WHERE u.first_name LIKE ? AND u.second_name LIKE ? ORDER BY u.id LIMIT 100",
		firstName+"%", lastName+"%")
	if err != nil {
		utils.Code500(c, "Query error", -7)
		return
	}
	defer rows.Close()

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
	if err = rows.Err(); err != nil {
		utils.Code500(c, "DB error", -9)
		return
	}

	c.JSON(200, users)
}
package routes

import (
	"github.com/gin-gonic/gin"
	"go-auth/controllers"
	"go-auth/middlewares"
)

func Route(r *gin.Engine) {

	r.Use(middlewares.RequestID())

	r.POST("/login", controllers.Login)
	r.POST("/user/register", controllers.UserRegister)
	r.GET("/user/get/:id", controllers.UserGetId)

	r.GET("/health", controllers.Health)

	authorized := r.Group("/test")
	authorized.Use(middlewares.AuthRequired())
	{
		authorized.GET("/test", controllers.AuthTest)
	}

}

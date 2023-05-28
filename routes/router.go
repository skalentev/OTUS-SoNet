package routes

import (
	"github.com/gin-gonic/gin"
	"otus-sonet/controllers"
	"otus-sonet/middlewares"
)

func Route(r *gin.Engine) {

	r.Use(middlewares.RequestID())

	r.POST("/login", controllers.Login)
	r.POST("/user/register", controllers.UserRegister)
	r.GET("/user/get/:id", controllers.UserGetId)
	r.GET("/user/search", controllers.UserSearch)

	r.GET("/health", controllers.Health)
	//	r.GET("/metrics", controllers.PrometheusHandler())
	authorized := r.Group("/test")
	authorized.Use(middlewares.AuthRequired())
	{
		authorized.GET("/test", controllers.AuthTest)
	}

	r.NoMethod(controllers.NoRoute)
	r.NoRoute(controllers.NoRoute)

}

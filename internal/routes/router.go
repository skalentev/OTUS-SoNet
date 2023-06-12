package routes

import (
	"github.com/gin-gonic/gin"
	controllers2 "otus-sonet/internal/controllers"
	middlewares2 "otus-sonet/internal/middlewares"
)

func Route(r *gin.Engine) {

	r.Use(middlewares2.RequestID())

	r.POST("/login", controllers2.Login)
	r.POST("/user/register", controllers2.UserRegister)
	r.GET("/user/get/:id", controllers2.UserGetId)
	r.GET("/user/search", controllers2.UserSearch)

	r.GET("/health", controllers2.Health)
	//	r.GET("/metrics", controllers.PrometheusHandler())
	authorized := r.Group("/test")
	authorized.Use(middlewares2.AuthRequired())
	{
		authorized.GET("/test", controllers2.AuthTest)
	}

	r.NoMethod(controllers2.NoRoute)
	r.NoRoute(controllers2.NoRoute)

}

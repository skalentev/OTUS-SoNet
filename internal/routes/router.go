package routes

import (
	"github.com/gin-gonic/gin"
	"otus-sonet/internal/controllers"
	"otus-sonet/internal/middlewares"
)

func Route(r *gin.Engine) {

	r.Use(middlewares.RequestID())

	r.POST("/login", controllers.Login)
	r.POST("/user/register", controllers.UserRegister)
	r.GET("/user/get/:id", controllers.UserGetId)
	r.GET("/user/search", controllers.UserSearch)
	r.GET("/post/get/:id", controllers.PostGet)
	r.GET("/friend/get/:user_id", controllers.FriendGet)
	r.GET("/health", controllers.Health)
	//	r.GET("/metrics", controllers.PrometheusHandler())
	authorized := r.Group("/")
	authorized.Use(middlewares.AuthRequired())
	{
		authorized.PUT("/friend/set/:user_id", controllers.FriendSet)
		authorized.DELETE("/friend/delete/:user_id", controllers.FriendDelete)
		authorized.PUT("/post/delete/:id", controllers.PostDelete)
		authorized.PUT("/post/update", controllers.PostUpdate)
		authorized.POST("/post/create", controllers.PostCreate)

		authorized.GET("/test", controllers.AuthTest)
	}

	r.NoMethod(controllers.NoRoute)
	r.NoRoute(controllers.NoRoute)

}

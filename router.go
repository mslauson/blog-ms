package main

import (
	"net/http"

	"gitea.slauson.io/slausonio/go-utils/siomw"
	"gitea.slauson.io/slausonio/iam-ms/controller"
	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(siomw.ErrorHandler)

	uc := controller.NewUserController()
	sc := controller.NewSessionController()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/iam", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// App
	v1 := r.Group("/api/iam/v1", siomw.AuthMiddleware)
	{
		user := v1.Group("/user")
		{
			user.GET("", uc.ListUsers)
			user.POST("", uc.CreateUser)
			user.GET("/:userId", uc.GetUserById)
			user.PUT("/:userId/password", uc.UpdatePassword)
			user.PUT("/:userId/email", uc.UpdateEmail)
			user.PUT("/:userId/phone", uc.UpdatePhone)
			user.DELETE("/:id", uc.DeleteUser)
		}

		session := v1.Group("/session")
		{
			session.POST("/email", sc.CreateEmailSession)
			session.DELETE("/:sessionId", sc.DeleteSession)
		}
	}

	return r
}

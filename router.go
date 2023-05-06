package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gitea.slauson.io/slausonio/go-utils/siomw"
	"gitea.slauson.io/slausonio/iam-ms/controller"
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

	v1 := r.Group("/api/iam/v1", siomw.AuthMiddleware)
	{
		user := v1.Group("/user")
		{
			user.GET("", uc.ListUsers)
			user.POST("", uc.CreateUser)
			user.GET("/:id", uc.GetUserById)
			user.PUT("/:id/password", uc.UpdatePassword)
			user.PUT("/:id/email", uc.UpdateEmail)
			user.PUT("/:id/phone", uc.UpdatePhone)
			user.DELETE("/:id", uc.DeleteUser)
		}

		session := v1.Group("/session")
		{
			session.POST("/email", sc.CreateEmailSession)
			session.DELETE("/:id/:sessionId", sc.DeleteSession)
		}
	}

	r.GET("/api/iam/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

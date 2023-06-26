package handler

import (
	"net/http"

	"gitea.slauson.io/slausonio/go-utils/siomw"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(siomw.PrometheusMiddleware())
	r.Use(siomw.ErrorHandler)

	h := NewBlogHdlr()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/post", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// App
	v1 := r.Group("/api/v1/post", siomw.AuthMiddleware)
	{
		v1.POST("", h.CreatePost)
		v1.GET("", h.GetAllPosts)
		id := v1.Group("/:id")
		{
			id.GET("", h.GetPost)
			id.PATCH("", h.UpdatePost)
			id.DELETE("", h.SoftDeletePost)
		}

		comment := v1.Group("/comment")
		{
			comment.POST("", h.AddComment)
			commentId := comment.Group("/:id")
			{
				commentId.PATCH("", h.UpdateComment)
				commentId.DELETE("", h.SoftDeleteComment)
			}
		}
	}

	r.GET("/api/blog/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

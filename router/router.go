package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netdisk/controller"
	"netdisk/middleware"
)

func Start() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	user := r.Group("/user", middleware.NeedNoLogin())
	{
		user.POST("/register", controller.Register)
		user.GET("/login", controller.Login)
	}

	file := r.Group("/file", middleware.JWT())
	{
		file.POST("/upload", controller.Upload)
		file.GET("/query", controller.Query)

		file.GET("/share",controller.Share)
		file.PUT("/change_path",controller.ChangePath)
		file.PUT("/rename",controller.Rename)
		file.DELETE("/delete",controller.DeleteFile)
	}
	r.GET("/file/download",controller.Download)


	r.Run(":8080")
}

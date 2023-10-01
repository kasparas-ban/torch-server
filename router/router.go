package router

import (
	"torch/torch-server/auth"
	"torch/torch-server/controllers"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter(logging, useAuth bool) *gin.Engine {
	if logging {
		Router = gin.Default()
	} else {
		Router = gin.New()
	}

	RegisterRoutes(Router, useAuth)

	return Router
}

func RegisterRoutes(r *gin.Engine, useAuth bool) *gin.Engine {
	r.Use(CORSMiddleware())
	if useAuth {
		r.Use(auth.AuthMiddleware())
	}

	api := r.Group("/api")
	{
		api.GET("/user-info", controllers.GetUserInfo)

		api.GET("/items", controllers.GetAllItems)
		api.POST("/add-item", controllers.AddItem)
		api.DELETE("/remove-item/:itemID", controllers.RemoveItem)
		api.PUT("/update-item", controllers.UpdateItem)
		api.PUT("/update-item-progress", controllers.UpdateItemProgress)

		api.GET("/timer-history", controllers.GetTimerHistory)
		api.PUT("/add-timer-record", controllers.UpsertTimerHistory)
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

package controllers

import (
	"torch/torch-server/auth"
	"torch/torch-server/controllers/items"
	"torch/torch-server/controllers/users"

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
		api.GET("/user-info", users.GetUserInfo)

		api.GET("/items", items.GetAllItems)
		api.POST("/add-item/:type", items.AddItem)

		api.DELETE("/remove-item/:itemID", items.RemoveItem)
		api.PUT("/update-item", items.UpdateItem)
		api.PUT("/update-item-progress", items.UpdateItemProgress)

		api.GET("/timer-history", GetTimerHistory)
		api.PUT("/add-timer-record", UpsertTimerHistory)
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

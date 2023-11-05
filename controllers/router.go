package controllers

import (
	"os"
	"torch/torch-server/auth"
	"torch/torch-server/controllers/history"
	"torch/torch-server/controllers/items"
	"torch/torch-server/controllers/users"
	"torch/torch-server/models"

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

	models.InitializeValidators()
	gin.EnableJsonDecoderDisallowUnknownFields()

	return Router
}

func RegisterRoutes(r *gin.Engine, useAuth bool) *gin.Engine {
	r.Use(CORSMiddleware())

	public := r.Group("/api")
	if useAuth {
		public.Use(auth.AuthMiddleware(true))
	}
	public.POST("/add-user", users.HandleAddNewUser)

	api := r.Group("/api")
	if useAuth {
		api.Use(auth.AuthMiddleware(false))
	}
	{
		api.GET("/user-info", users.HandleGetUserInfo)
		api.PUT("/update-user", users.HandleUpdateUser)
		api.DELETE("/delete-user", users.HandleDeleteUser)

		api.GET("/items", items.HandleGetAllItems)
		api.GET("/item/:itemID", items.HandleGetItem)
		api.POST("/add-item/:type", items.HandleAddItem)
		api.DELETE("/remove-item/:itemID", items.HandleRemoveItem)
		api.PUT("/update-item/:type", items.HandleUpdateItem)
		api.PUT("/update-item-progress", items.HandleUpdateItemProgress)

		api.GET("/timer-history", history.HandleGetTimerHistory)
		api.PUT("/add-timer-record", history.HandleUpsertTimerHistory)
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", getOrigin(c))
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

func getOrigin(c *gin.Context) string {
	allowedOrigin := []string{os.Getenv("LOCAL_FE_DOMAIN"), os.Getenv("DEV_FE_DOMAIN"), os.Getenv("PROD_FE_DOMAIN")}

	origin := c.GetHeader("Origin")
	returnOrigin := allowedOrigin[0]

	for _, element := range allowedOrigin {
		if element == origin {
			returnOrigin = origin
			break
		}
	}

	return returnOrigin
}

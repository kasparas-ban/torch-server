package controllers

import (
	"net/http"
	"os"
	"torch/torch-server/auth"
	"torch/torch-server/controllers/history"
	"torch/torch-server/controllers/items"
	"torch/torch-server/controllers/notify"
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

	r.GET("/ping", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	public := r.Group("/api")
	{
		public.POST("/add-user", users.HandleAddNewUser)
		public.POST("/notify", notify.HandleNotifyEmail)
	}

	private := r.Group("/api")
	if useAuth {
		private.Use(auth.AuthMiddleware(false))
	}
	{
		private.GET("/user-info", users.HandleGetUserInfo)
		private.PUT("/update-user", users.HandleUpdateUser)
		private.PUT("/update-user-email", users.HandleUpdateUserEmail)
		private.DELETE("/delete-user", users.HandleDeleteUser)

		private.GET("/items", items.HandleGetAllItems)
		private.GET("/item/:itemID", items.HandleGetItem)
		private.POST("/add-item/:type", items.HandleAddItem)
		private.DELETE("/delete-item", items.HandleRemoveItem)
		private.PUT("/update-item/:type", items.HandleUpdateItem)
		private.PUT("/update-status", items.HandleUpdateItemStatus)
		private.PUT("/update-item-progress", items.HandleUpdateItemProgress)
		private.PUT("/update-user-progress", items.HandleUpdateUserProgress)

		private.GET("/timer-history", history.HandleGetTimerHistory)
		private.PUT("/add-timer-record", history.HandleUpsertTimerHistory)
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

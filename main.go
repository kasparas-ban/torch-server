package main

import (
	"fmt"
	"log"
	"os"

	"torch/torch-server/auth"
	c "torch/torch-server/controllers"
	"torch/torch-server/db"

	uuid "github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	db.Init()
	auth.Init()

	r := gin.Default()

	//Custom form validator
	// binding.Validator = new(forms.DefaultValidator)

	r.Use(CORSMiddleware())
	// r.Use(auth.AuthMiddleware())
	// r.Use(RequestIDMiddleware())

	api := r.Group("/api")
	{
		api.GET("/full-profile/:clerkID", c.GetUserInfoByClerkID)
		api.GET("/user/:userID", c.GetUserInfo)

		api.GET("/items/:userID", c.GetAllItems)
		api.POST("/add-item/:userID", c.AddItem)
		api.PUT("/update-item/:userID", c.UpdateItem)


		/*** START Article ***/
		// article := new(controllers.ArticleController)

		// v1.POST("/article", TokenAuthMiddleware(), article.Create)
		// v1.GET("/articles", TokenAuthMiddleware(), article.All)
		// v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
		// v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
		// v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)
	}

	// r.NoRoute(func(c *gin.Context) {
	// 	c.HTML(404, "404.html", gin.H{})
	// })

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {
		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}
}

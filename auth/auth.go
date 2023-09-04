package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

var client clerk.Client

func Init() {
	var err error
	client, err = clerk.NewClient(os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := c.GetHeader("Authorization")
		sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")

		_, err := client.VerifyToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetClerkClient() clerk.Client {
	return client
}

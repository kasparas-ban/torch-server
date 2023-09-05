package auth

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
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

		sessClaims, err := client.VerifyToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user, err := client.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Failed to get user info")
			c.Abort()
			return
		}

		c.Set("userID", user.ID)
		c.Next()
	}
}

func GetClerkClient() clerk.Client {
	return client
}

func GetUserID(c *gin.Context) uint64 {
	userIDString := c.GetString("userID")
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Could not find user info")},
		)
		c.Abort()
		return 0
	}
	return userID
}
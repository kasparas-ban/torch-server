package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"torch/torch-server/models"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

var client clerk.Client

const (
	userID_metadata = "user_id"
	userID_context = "userID"
)

func Init() {
	var err error
	client, err = clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify session token
		sessionToken := c.GetHeader("Authorization")
		sessionToken = strings.TrimPrefix(sessionToken, "Bearer ")
		sessClaims, err := client.VerifyToken(sessionToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Read clerkID
		user, err := client.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": errors.New("Could not get user info")},
			)
			c.Abort()
			return
		}

		c.Set("clerkID", user.ID)

		// Read userID
		userIDString := user.PrivateMetadata.(map[string]interface{})[userID_metadata]
		userID, err := strconv.ParseUint(userIDString.(string), 10, 64)
		if err == nil {
			c.Set(userID_context, userID)
		} else {
			err := addUserID(c, user)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": errors.New("Unexpected error occured")},
				)
				c.Abort()
			}
		}

		c.Next()
	}
}

func GetClerkClient() clerk.Client {
	return client
}

func GetClerkID(c *gin.Context) string {
	clerkID := c.GetString("clerkID")
	if clerkID != "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": errors.New("Could not find clerkID")},
		)
		c.Abort()
	}
	return clerkID
}

func GetUserID(c *gin.Context) uint64 {
	userID := c.GetUint64(userID_context)
	if userID == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Could not find userID"},
		)
		c.Abort()
		return 0
	}
	return userID
}

func addUserID(c *gin.Context, user *clerk.User) error {	
	userInfo, err := models.GetUserInfoByClerkID(user.ID)
	if err != nil {
		return err
	}

	currPrivMetadata := make(map[string]interface{})
	currPrivMetadata[userID_metadata] = strconv.FormatUint(userInfo.UserID, 10)

	var keyValStrings []string
	for key, value := range currPrivMetadata {
		keyValStrings = append(keyValStrings, fmt.Sprintf("\"%s\": \"%s\"", key, value))
	}
	finalStr := fmt.Sprintf("{%s}", strings.Join(keyValStrings, ","))
	_, err = client.Users().Update(user.ID, &clerk.UpdateUser{
		PrivateMetadata: finalStr,
	})
	if err != nil {
		return err
	}

	c.Set(userID_context, userInfo.UserID)

	return nil
}
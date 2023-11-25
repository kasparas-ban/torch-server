package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"torch/torch-server/models"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

var client clerk.Client

const (
	userID_metadata = "user_id"
	userID_context  = "userID"
)

func Init() {
	var err error
	client, err = clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	clerk.WithSessionV2(
		client,
		clerk.WithLeeway(5*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func AuthMiddleware(isNewUser bool) gin.HandlerFunc {
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
				gin.H{"error": "Could not get user info"},
			)
			c.Abort()
			return
		}

		c.Set("clerkID", user.ID)
		c.Set("setClerkMetadata", func() error {
			return addUserID(c, user)
		})

		if !isNewUser {
			// Read userID
			userIDString := user.PrivateMetadata.(map[string]interface{})[userID_metadata]
			if userIDString == nil {
				fmt.Printf("\n\n METADATA NULL \n\n")
				err := addUserID(c, user)
				if err != nil {
					c.JSON(
						http.StatusInternalServerError,
						gin.H{"error": "Unexpected error occured"},
					)
					c.Abort()
				}
			} else {
				userID, err := strconv.ParseUint(userIDString.(string), 10, 64)
				fmt.Printf("\n\n METADATA EXISTS %v \n\n", userID)
				if err == nil {
					c.Set(userID_context, userID)
				} else {
					err := addUserID(c, user)
					if err != nil {
						c.JSON(
							http.StatusInternalServerError,
							gin.H{"error": "Unexpected error occured"},
						)
						c.Abort()
					}
				}
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
			gin.H{"error": "Could not find clerkID"},
		)
		c.Abort()
	}
	return clerkID
}

func GetUserID(c *gin.Context) (uint64, error) {
	userID := c.GetUint64(userID_context)
	if userID == 0 {
		return 0, errors.New("Failed to get user ID")
	}
	return userID, nil
}

func addUserID(c *gin.Context, user *clerk.User) error {
	userInfo, err := models.GetUserByClerkID(user.ID)
	if err != nil {
		return err
	}
	if userInfo.UserID == 0 {
		return errors.New("Failed to get User ID")
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

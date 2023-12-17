package users

import (
	"encoding/json"
	"net/http"
	"torch/torch-server/models"

	"github.com/gin-gonic/gin"
)

type EmailAddresses struct {
	EmailAddress string `json:"email_address"`
}

type NewReqData struct {
	ID             string           `json:"id"`
	Username       string           `json:"username"`
	EmailAddresses []EmailAddresses `gorm:"embedded" json:"email_addresses"`
}

type NewUserReq struct {
	Data NewReqData `gorm:"embedded" json:"data"`
}

func HandleAddNewUser(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)

	var data NewUserReq
	err := decoder.Decode(&data)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	if data.Data.ID == "" || len(data.Data.EmailAddresses) == 0 {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid user object"},
		)
		return
	}

	newUser := models.NewUser{
		ClerkID: data.Data.ID,
		Email:   data.Data.EmailAddresses[0].EmailAddress,
	}

	user, err := models.AddUser(newUser)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	setClerkMetadata, exists := c.Get("setClerkMetadata")
	setClerkFunc, ok := setClerkMetadata.(func() error)
	if !ok || !exists || setClerkFunc() != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Unexpected error occured"},
		)
		c.Abort()
	}

	c.JSON(http.StatusOK, user)
}

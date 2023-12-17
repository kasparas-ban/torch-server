package notify

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type NotifyReq struct {
	EmailAddress string
}

func HandleNotifyEmail(c *gin.Context) {
	var emailData NotifyReq
	if err := c.BindJSON(&emailData); err != nil {
		fmt.Printf("Error: %v", err)
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid email"},
		)
		return
	}

	email := os.Getenv("EMAIL")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	hostAddress := os.Getenv("EMAIL_HOST_ADDRESS")
	hostPortStr := os.Getenv("EMAIL_PORT")
	hostPort, err := strconv.Atoi(hostPortStr)
	if err != nil {
		log.Fatal(err)
	}

	emailBody := fmt.Sprintf("Notify me when new features get added.<br><br>Email address: %v", emailData.EmailAddress)

	// Form the email
	m := gomail.NewMessage()
	m.SetHeader("From", "Torch App <"+email+">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Torch: Notify for updates")
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		hostAddress,
		hostPort,
		email,
		emailPassword,
	)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "internal server error"},
		)
		return
	}

	c.JSON(http.StatusOK, nil)
}

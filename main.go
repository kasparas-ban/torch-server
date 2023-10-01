package main

import (
	"flag"
	"log"
	"os"

	"torch/torch-server/auth"
	"torch/torch-server/db"
	"torch/torch-server/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	port := os.Getenv("PORT")
	prod := flag.Bool("prod", false, "Production environment")
	flag.Parse()

	if *prod {
		gin.SetMode(gin.ReleaseMode)
	}

	db.Init(os.Getenv("DSN"))
	auth.Init()

	r := router.SetupRouter(true, true)

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {
		// Generated using sh generate-certificate.sh
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

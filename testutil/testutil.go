package testutil

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"torch/torch-server/db"
	"torch/torch-server/router"

	"github.com/gin-gonic/gin"
)

var TestRouter *gin.Engine
var MockUser uint64

func TestMain(m *testing.M) {
	dbUsername := "root"
	dbPassword := "root"
	dbName := "torch-database"
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?parseTime=true", dbUsername, dbPassword, dbName)
	db.Init(dsn)

	gin.SetMode(gin.TestMode)
	TestRouter = router.SetupRouter(false, false)

	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware(r *gin.Engine) {
	MockUser = 1
	r.Use(func(c *gin.Context) {
		c.Set("userID", MockUser)
	})
}

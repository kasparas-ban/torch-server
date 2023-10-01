package testutil

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"torch/torch-server/db"

	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

var (
	dbUsername = "root"
	dbPassword = "root"
	dbName     = "torch-database"
)

var MockUser uint64

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

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

func PrepareMySQLContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUsername),
		mysql.WithPassword(dbPassword),
		mysql.WithScripts("../db/init.sql"),
	)
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port("3306"))
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, hostIP, mappedPort.Port(), dbName)
	db.Init(dsn)

	log.Printf("TestContainers: container %s is now running at %s\n", "mysql:latest", hostIP)
	return container, nil
}

func ContainerCleanUp(ctx context.Context, container *mysql.MySQLContainer) {
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
}

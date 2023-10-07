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
	DBUsername = "root"
	DBPassword = "root"
	DBName     = "torch-database"
)

var MockUser uint64
var Ctx context.Context
var MysqlContainer *mysql.MySQLContainer

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	local := flag.Bool("local", false, "Using locally running MySQL database")
	flag.Parse()

	if *local {
		// Making a local database connection
		dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?multiStatements=true", DBUsername, DBPassword, DBName)
		db.Init(dsn)
	} else {
		// MySQL database container setup
		Ctx = context.Background()
		container, err := PrepareMySQLContainer(Ctx)
		if err != nil {
			panic(err)
		}
		defer ContainerCleanUp(Ctx, container)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func MockAuthMiddleware(r *gin.Engine, mockUserID uint64) {
	MockUser = mockUserID
	r.Use(func(c *gin.Context) {
		c.Set("userID", MockUser)
	})
}

func PrepareMySQLContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:latest"),
		mysql.WithDatabase(DBName),
		mysql.WithUsername(DBUsername),
		mysql.WithPassword(DBPassword),
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", DBUsername, DBPassword, hostIP, mappedPort.Port(), DBName)
	db.Init(dsn)

	log.Printf("TestContainers: container %s is now running at %s\n", "mysql:latest", hostIP)
	return container, nil
}

func ContainerCleanUp(ctx context.Context, container *mysql.MySQLContainer) {
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
}

func getSQL(filename string) (string, error) {
	// This will run from /tests directory
	filePath := fmt.Sprintf("../testutil/sqlScripts/%s.sql", filename)
	sqlContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(sqlContents), nil
}

func SeedDB() {
	// Users data
	filename := "insertUsers"
	cmd, err := getSQL(filename)
	if err != nil {
		panic("Failed to load SQL command")
	}

	err = db.GetDB().Exec(cmd).Error
	if err != nil {
		panic("Failed to seed the DB")
	}

	// Items data
	filename = "insertItems"
	cmd, err = getSQL(filename)
	if err != nil {
		panic("Failed to load SQL command")
	}

	err = db.GetDB().Exec(cmd).Error
	if err != nil {
		panic("Failed to seed the DB")
	}
}

func CleanAllTables() {
	filename := "cleanAllTables"
	cmd, err := getSQL(filename)
	if err != nil {
		panic("Failed to load SQL command")
	}

	err = db.GetDB().Exec(cmd).Error
	if err != nil {
		panic("Failed to clean all tables")
	}
}

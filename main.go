package main

import (
	"fullstack2024-test/database"
	"fullstack2024-test/handler"
	"fullstack2024-test/repository"
	"fullstack2024-test/service"
	"fullstack2024-test/usecase"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(database.InitDB()), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database!")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to connect to database!")
	}
	defer sqlDB.Close()

	database.InitRedis()

	s3Service, err := service.NewS3Service()
	if err != nil {
		log.Printf("Warning: Failed to initialize AWS S3 service: %v", err)
		s3Service = nil
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	clientRepo := repository.NewClientRepository(db)
	clientUseCase := usecase.NewClientUseCase(clientRepo)
	clientHandler := handler.NewClientHandler(clientUseCase, s3Service)
	clientHandler.ClientRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

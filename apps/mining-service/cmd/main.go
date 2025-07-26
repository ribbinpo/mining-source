package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ribbinpo/mining-service/internal/application/usecase"
	"github.com/ribbinpo/mining-service/internal/config"
	"github.com/ribbinpo/mining-service/internal/frameworks/database"
	repository "github.com/ribbinpo/mining-service/internal/frameworks/database/repository"
	"github.com/ribbinpo/mining-service/internal/frameworks/webscraping"
	"github.com/ribbinpo/mining-service/internal/frameworks/webservice"
	"github.com/ribbinpo/mining-service/internal/frameworks/webservice/router"
)

func main() {
	config := config.NewConfig(".env")
	webService := webservice.NewWebService(config.App.Port)

	api := webService.FiberApp.Group("/api")

	databaseInstance := database.NewDatabase(config.Db.Dsn)
	db := databaseInstance.Connect()
	scraperRepository := webscraping.NewScraperRepository()
	stoneRepository := repository.NewStoneRepository(db)
	stoneUsecase := usecase.NewStoneUsecase(scraperRepository, stoneRepository)

	router.StoneRouter(api, stoneUsecase)

	webService.FiberApp.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Mining service is running")
		return c.SendString("Mining service is running")
	})

	webService.Start()
}

package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ribbinpo/mining-service/internal/application/port"
	"github.com/ribbinpo/mining-service/internal/frameworks/webservice/handler"
)

func StoneRouter(router fiber.Router, stoneUsecase port.StoneUsecase) fiber.Router {
	stoneRouter := router.Group("/stones")
	stoneHandler := handler.NewStoneHandler(stoneUsecase)

	stoneRouter.Get("/:url", stoneHandler.GetStoneByURL)
	stoneRouter.Post("/", stoneHandler.RegisterDig)

	return stoneRouter
}

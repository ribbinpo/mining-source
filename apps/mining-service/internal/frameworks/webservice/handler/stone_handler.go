package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ribbinpo/mining-service/internal/application/port"
	"github.com/ribbinpo/mining-service/internal/util"
)

type StoneHandler struct {
	StoneUsecase port.StoneUsecase
}

func NewStoneHandler(stoneUsecase port.StoneUsecase) *StoneHandler {
	return &StoneHandler{
		StoneUsecase: stoneUsecase,
	}
}

func (h *StoneHandler) RegisterDig(c *fiber.Ctx) error {
	type Dig struct {
		Url     string
		Refresh bool
	}
	d := new(Dig)
	if err := c.BodyParser(d); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := h.StoneUsecase.RegisterDig(d.Url, d.Refresh); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Dig registered successfully",
	})
}

func (h *StoneHandler) GetStoneByURL(c *fiber.Ctx) error {
	encodedURL := c.Params("url")

	// Decode the URL parameter
	decodedURL, err := util.DecodeURL(encodedURL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL encoding: " + err.Error(),
		})
	}

	stone, err := h.StoneUsecase.GetStoneByURL(decodedURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(stone)
}

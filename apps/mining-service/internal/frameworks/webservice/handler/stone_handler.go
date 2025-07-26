package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ribbinpo/mining-service/internal/application/port"
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
	stone, err := h.StoneUsecase.GetStoneByURL(c.Params("url"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(stone)
}

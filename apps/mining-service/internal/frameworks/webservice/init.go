package webservice

import "github.com/gofiber/fiber/v2"

type WebService struct {
	FiberApp *fiber.App
	Port     string
}

func NewWebService(port string) *WebService {
	fiberApp := fiber.New()
	return &WebService{
		FiberApp: fiberApp,
		Port:     port,
	}
}

func (w *WebService) Start() {
	w.FiberApp.Listen(w.Port)
}

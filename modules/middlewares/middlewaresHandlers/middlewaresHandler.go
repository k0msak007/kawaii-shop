package middlewaresHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/k0msak007/kawaii-shop/config"
	"github.com/k0msak007/kawaii-shop/modules/entities"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresUsecases"
)

type middlewaresHandlersError string

const (
	routerCheckErr middlewaresHandlersError = "router-001"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type middlewaresHandler struct {
	cfg                 config.IConfig
	middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases) IMiddlewaresHandler {
	return &middlewaresHandler{
		middlewaresUsecases: middlewaresUsecases,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:          cors.ConfigDefault.Next,
		AllowOrigins:  "*",
		AllowHeaders:  "",
		AllowMethods:  " POST, HEAD, PUT, DELETE, PATCH",
		ExposeHeaders: "",
		MaxAge:        0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path} \n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Asia/Bangkok",
	})
}

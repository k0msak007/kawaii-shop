package middlewaresHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/k0msak007/kawaii-shop/config"
	"github.com/k0msak007/kawaii-shop/modules/entities"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/k0msak007/kawaii-shop/pkg/kawaiiauth"
)

type middlewaresHandlersError string

const (
	routerCheckErr middlewaresHandlersError = "router-001"
	jwtAuthErr     middlewaresHandlersError = "router-002"
	paramsCheckErr middlewaresHandlersError = "router-003"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
}

type middlewaresHandler struct {
	cfg                 config.IConfig
	middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                 cfg,
		middlewaresUsecases: middlewaresUsecases,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:          cors.ConfigDefault.Next,
		AllowOrigins:  "*",
		AllowHeaders:  "",
		AllowMethods:  "GET, POST, HEAD, PUT, DELETE, PATCH",
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

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := kawaiiauth.ParseToken(h.cfg.Jwt(), token)

		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecases.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(string)
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}

		return c.Next()
	}
}

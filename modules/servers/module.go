package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresHandlers"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresRepositories"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/k0msak007/kawaii-shop/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecases := middlewaresUsecases.MiddlewaresUsecases(repository)

	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecases)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

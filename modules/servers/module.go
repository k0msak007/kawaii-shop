package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresHandlers"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresRepositories"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/k0msak007/kawaii-shop/modules/monitor/monitorHandlers"
	"github.com/k0msak007/kawaii-shop/modules/users/usersHandlers"
	"github.com/k0msak007/kawaii-shop/modules/users/usersRepositories"
	"github.com/k0msak007/kawaii-shop/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecases := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecases)

	router := m.r.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefressPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)

	router.Get("/secret", handler.GenerateAdminToken)
}

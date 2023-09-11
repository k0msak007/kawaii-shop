package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoHandlers"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoRepositories"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoUsecases"
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
	AppinfoModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefressPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecases := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecases)

	router := m.r.Group("/appinfo")

	router.Get("/categorires", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
}

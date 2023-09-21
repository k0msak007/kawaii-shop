package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoHandlers"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoRepositories"
	"github.com/k0msak007/kawaii-shop/modules/appinfo/appinfoUsecases"
	"github.com/k0msak007/kawaii-shop/modules/files/filesHandlers"
	"github.com/k0msak007/kawaii-shop/modules/files/filesUsecases"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresHandlers"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresRepositories"
	"github.com/k0msak007/kawaii-shop/modules/middlewares/middlewaresUsecases"
	"github.com/k0msak007/kawaii-shop/modules/monitor/monitorHandlers"
	"github.com/k0msak007/kawaii-shop/modules/orders/ordersHandlers"
	"github.com/k0msak007/kawaii-shop/modules/orders/ordersRepositories"
	"github.com/k0msak007/kawaii-shop/modules/orders/ordersUsecases"
	"github.com/k0msak007/kawaii-shop/modules/products/productsHandlers"
	"github.com/k0msak007/kawaii-shop/modules/products/productsRepositories"
	"github.com/k0msak007/kawaii-shop/modules/products/productsUsecases"
	"github.com/k0msak007/kawaii-shop/modules/users/usersHandlers"
	"github.com/k0msak007/kawaii-shop/modules/users/usersRepositories"
	"github.com/k0msak007/kawaii-shop/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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
	router.Post("/signin", handler.SignIn)
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

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecases := filesUsecases.FileUsecase(m.s.cfg)
	handler := filesHandlers.FileHandler(m.s.cfg, usecases)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	filesUsecases := filesUsecases.FileUsecase(m.s.cfg)

	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecases)
	productsUsecases := productsUsecases.ProductsUsecase(productsRepository)
	productsHandler := productsHandlers.ProductsHandler(m.s.cfg, productsUsecases, filesUsecases)

	router := m.r.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.AddProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.UpdateProduct)

	router.Get("/", m.mid.ApiKeyAuth(), productsHandler.FindProduct)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), productsHandler.FindOneProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.DeleteProduct)
}

func (m *moduleFactory) OrdersModule() {
	filesUsecases := filesUsecases.FileUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecases)

	ordersRepository := ordersRepositories.OrderRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrderUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrderHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), ordersHandler.FindOneOrder)

	router.Patch("/:user_id/:order_id", m.mid.JwtAuth(), ordersHandler.UpdateOrder)
}

package servers

import (
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoHandlers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoRepositories"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/appinfo/appinfoUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files/filesHandlers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/files/filesUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareHandlers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareRepositories"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/middlewares/middlewareUsecases"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/monitor/monitorHandlers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/users/usersHandlers"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/users/usersRepositories"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareHandlers.IMiddelwaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewareHandlers.IMiddelwaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewareHandlers.IMiddelwaresHandler {
	repository := middlewareRepositories.MiddlewaresRespository(s.db)
	usecase := middlewareUsecases.MiddlewaresUsecase(repository)
	return middlewareHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	// /v1/users/sign

	router := m.r.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignOut)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecases(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}

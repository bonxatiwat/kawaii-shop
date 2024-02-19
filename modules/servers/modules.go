package servers

import (
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

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
}

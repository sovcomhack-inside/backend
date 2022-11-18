package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sovcomhack-inside/internal/api/controller"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} else if e, ok := err.(*constants.CodedError); ok {
		code = e.Code()
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return ctx.Status(code).SendString(err.Error())
}

type APIService struct {
	log    *logger.Logger
	router *fiber.App
}

func (svc *APIService) Serve(addr string) {
	svc.log.Fatal(svc.router.Listen(addr))
}

func (svc *APIService) Shutdown() error {
	return svc.router.Shutdown()
}

func NewAPIService(log *logger.Logger, dbRegistry store.Store) (*APIService, error) {
	svc := &APIService{
		log: log,
		router: fiber.New(fiber.Config{
			ErrorHandler: errorHandler,
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
		}),
	}

	service := service.NewService(log, dbRegistry)
	controller := controller.NewController(log, service)

	api := svc.router.Group("/api", recover.New())

	auth := api.Group("/auth")
	auth.Post("/signup", controller.SignupUser)
	auth.Post("/login", controller.LoginUser)
	auth.Delete("/logout", controller.LogoutUser)

	oauth := api.Group("/oauth")
	oauth.Get("/telegram", controller.OAuthTelegram)

	return svc, nil
}

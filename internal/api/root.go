package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sovcomhack-inside/internal/api/controller"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
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

func NewAPIService(log *logger.Logger, dbRegistry *store.Registry) (*APIService, error) {
	svc := &APIService{
		log: log,
		router: fiber.New(fiber.Config{
			ErrorHandler: errorHandler,
			JSONEncoder:  sonic.Marshal,
			JSONDecoder:  sonic.Unmarshal,
		}),
	}

	registry, err := controller.NewRegistry(log, dbRegistry)
	if err != nil {
		return nil, err
	}

	api := svc.router.Group("/api", recover.New())

	auth := api.Group("/auth")
	auth.Post("/signup", registry.AuthController.SignupUser)
	auth.Post("/login", registry.AuthController.LoginUser)
	auth.Post("/logout", registry.AuthController.LogoutUser)

	return svc, nil
}

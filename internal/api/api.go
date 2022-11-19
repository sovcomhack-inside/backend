package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sovcomhack-inside/internal/api/controller"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type APIService struct {
	router *echo.Echo
}

func (svc *APIService) Serve(addr string) {
	logger.Fatal(context.Background(), svc.router.Start(addr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	return svc.router.Shutdown(ctx)
}

func NewAPIService(store store.Store) (*APIService, error) {
	svc := &APIService{router: echo.New()}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()

	service := service.NewService(store)
	controller := controller.NewController(service)

	// svc.router.Use(svc.XRequestIDMiddleware(), svc.LoggingMiddleware(), svc.AccessLogMiddleware())

	api := svc.router.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/signup", controller.SignupUser)
	auth.POST("/login", controller.LoginUser)
	auth.DELETE("/logout", controller.LogoutUser)

	account := api.Group("/accounts")

	account.POST("/", controller.CreateAccount)

	oauth := api.Group("/oauth")
	oauth.GET("/telegram", controller.OAuthTelegram)

	admin := api.Group("/admin")
	admin.PUT("/update_user_status", controller.UpdateUserStatus)

	return svc, nil
}

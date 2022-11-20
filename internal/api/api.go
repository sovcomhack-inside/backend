package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sovcomhack-inside/internal/api/controller"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/service"
	"github.com/sovcomhack-inside/internal/pkg/store"
)

type APIService struct {
	router  *echo.Echo
	store   store.Store
	service service.Service
}

func (svc *APIService) Serve(addr string) {
	logger.Fatal(context.Background(), svc.router.Start(addr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	return svc.router.Shutdown(ctx)
}

func NewAPIService(store store.Store) (*APIService, error) {
	svc := &APIService{router: echo.New(), store: store}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()
	svc.router.Use(middleware.Logger())
	svc.router.HTTPErrorHandler = httpErrorHandler

	svc.service = service.NewService(store)

	api := svc.router.Group("/api/v1")
	controller := controller.NewController(store, svc.service)

	auth := api.Group("/auth")
	auth.POST("/signup", controller.SignupUser)
	auth.POST("/login", controller.LoginUser)
	auth.DELETE("/logout", controller.LogoutUser)

	user := api.Group("/user", svc.AuthMiddleware)
	user.GET("/get", controller.GetUser)

	account := api.Group("/accounts", svc.AuthMiddleware)

	account.POST("/create", controller.CreateAccount)
	account.GET("/list", controller.ListUserAccounts)
	account.POST("/refill", controller.RefillAccount)
	account.POST("/withdraw", controller.WithdrawFromAccount)
	account.POST("/buy", controller.MakePurchase)
	account.POST("/sell", controller.MakeSale)

	oauth := api.Group("/oauth")
	oauth.GET("/telegram", controller.OAuthTelegram, svc.OAuthTelegramMiddleware)

	admin := api.Group("/admin")
	admin.POST("/login", controller.LoginAdmin)
	admin.POST("/update_user_status", controller.UpdateUserStatus, svc.AdminMiddleware)
	admin.POST("/list_users", controller.ListUsers, svc.AdminMiddleware)

	currencies := api.Group("/currencies", svc.AuthMiddleware)
	currencies.GET("/list", controller.ListCurrencies)
	currencies.GET("/data", controller.GetCurrencyData)

	operations := api.Group("/operations", svc.AuthMiddleware)
	operations.POST("/list", controller.ListOperations)

	funtik := api.Group("/funtik")
	funtik.POST("/subscribe", controller.SubscribeToFuntik)

	return svc, nil
}

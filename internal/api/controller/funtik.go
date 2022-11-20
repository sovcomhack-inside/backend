package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/logger"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

func (c *Controller) SubscribeToFuntik(ctx echo.Context) error {
	request := &dto.SubscribeToFuntikRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}
	if ctx.Get(constants.CtxKeyUserID) == nil {
		return fmt.Errorf("can't resolve user_id")
	}
	request.UserID = ctx.Get(constants.CtxKeyUserID).(int64)

	acc, err := c.store.SubscribeToFuntik(ctx.Request().Context(), request)
	if err != nil {
		logger.Errorf(ctx.Request().Context(), "BLYAT")
		return err
	}

	funtikRUBAcc := core.Account{
		Number:   uuid.New(),
		UserID:   request.UserID,
		Currency: "RUB",
		Balance:  decimal.NewFromInt(0),
		ForBot:   true,
	}
	err = c.store.CreateAccount(ctx.Request().Context(), &funtikRUBAcc)
	if err != nil {
		return err
	}
	funtikUSDAcc := core.Account{
		Number:   uuid.New(),
		UserID:   request.UserID,
		Currency: "USD",
		Balance:  decimal.NewFromInt(0),
		ForBot:   true,
	}
	err = c.store.CreateAccount(ctx.Request().Context(), &funtikUSDAcc)
	if err != nil {
		return err
	}
	user, err := c.store.GetUserByID(ctx.Request().Context(), request.UserID)
	// запускаем фонового бота
	go runFuntik(ctx.Request().Context(), user, &funtikRUBAcc, &funtikUSDAcc)
	// да немного поломал архитектуру контроллер-сервис-репо, но если честно я немного не успеваю, да и Оганес
	// сам ее уже сломал до меня!!
	logger.Error(ctx.Request().Context(), acc)
	operations := []*core.Operation{
		{
			Purpose:           "Месячная подписка на бота Фунтика",
			OperationType:     core.OperationTypeSubscription,
			AccountNumberFrom: &request.AccountNumberFrom,
			AmountCentsFrom:   lo.ToPtr(decimal.NewFromInt(request.SubscribePriceCents)),
			CurrencyFrom:      lo.ToPtr(acc.Currency),
			ExchangeRateRatio: 1,
		},
	}
	err = c.store.InsertOperations(ctx.Request().Context(), operations)
	if err != nil {
		return err
	}

	return nil
}

func runFuntik(ctx context.Context, user *core.User, funtikRUBAcc, funtikUSDAcc *core.Account) {
	for {
		if user.SubscriptionExpiredAt != nil && user.SubscriptionExpiredAt.After(time.Now()) {
			break
		}
		logger.Errorf(ctx, "funtik does something")
		// do some work
		time.Sleep(1 * time.Second)
	}
}

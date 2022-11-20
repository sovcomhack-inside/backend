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
		Balance:  decimal.NewFromInt(300000),
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
	go c.runFuntik(ctx.Request().Context(), user, &funtikRUBAcc, &funtikUSDAcc)
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
			ExchangeRateRatio: decimal.NewFromInt(1),
		},
	}
	err = c.store.InsertOperations(ctx.Request().Context(), operations)
	if err != nil {
		return err
	}

	return nil
}

type funtikStuff struct {
	IsNextOperationBuy   bool
	UpwardTrendThreshold float64
	DipThreshold         float64
	ProfitThreshold      float64
	StopLossThreshold    float64
	LastOpPrice          decimal.Decimal
	CurrentPrice         decimal.Decimal

	accountRUB *core.Account
	accountUSD *core.Account
}

func (c *Controller) runFuntik(ctx context.Context, user *core.User, funtikRUBAcc, funtikUSDAcc *core.Account) {
	stuff := funtikStuff{
		IsNextOperationBuy:   true,
		UpwardTrendThreshold: 1.50,
		DipThreshold:         -2.25,
		ProfitThreshold:      1.25,
		StopLossThreshold:    -2.00,

		accountRUB: funtikRUBAcc,
		accountUSD: funtikUSDAcc,
	}
	for {
		if user.SubscriptionExpiredAt != nil && user.SubscriptionExpiredAt.After(time.Now()) {
			break
		}
		err := c.decideToBuy(ctx, &stuff)
		if err != nil {
			logger.Errorf(ctx, "error in funtik occures: %w", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) decideToBuy(ctx context.Context, staff *funtikStuff) error {
	var err error
	var resp *dto.MakePurchaseResponse
	var percentageDiff float64

	percentageDiff = (staff.CurrentPrice.Sub(staff.LastOpPrice)).Div(staff.LastOpPrice).InexactFloat64() * float64(100)

	if staff.IsNextOperationBuy {
		// try to buy
		if percentageDiff <= staff.UpwardTrendThreshold || percentageDiff <= staff.DipThreshold {
			staff.CurrentPrice, err = c.service.LatestCurrencyPrice(ctx, staff.accountUSD.Currency, staff.accountRUB.Currency)
			resp, err = c.service.MakeTransfer(ctx, &dto.TransferRequestDTO{
				AccountFrom:       staff.accountRUB.Number,
				CreditAmountCents: staff.accountRUB.Balance.Mul(decimal.NewFromFloat(0.1)),
				AccountTo:         staff.accountUSD.Number,
				DebitAmountCents:  staff.accountRUB.Balance.Mul(decimal.NewFromFloat(0.1)).Mul(staff.CurrentPrice),
			}, staff.CurrentPrice)
			if err != nil {
				return err
			}
			staff.accountRUB = resp.NewAccountFrom
			staff.accountUSD = resp.NewAccountTo
			staff.LastOpPrice = staff.CurrentPrice
			staff.IsNextOperationBuy = false
		}
	} else {
		// try to sell
		staff.CurrentPrice, err = c.service.LatestCurrencyPrice(ctx, staff.accountRUB.Currency, staff.accountUSD.Currency)
		if percentageDiff >= staff.ProfitThreshold || percentageDiff <= staff.StopLossThreshold {
			resp, err = c.service.MakeTransfer(ctx, &dto.TransferRequestDTO{
				AccountFrom:       staff.accountUSD.Number,
				CreditAmountCents: staff.accountUSD.Balance.Mul(decimal.NewFromFloat(0.1)),
				AccountTo:         staff.accountRUB.Number,
				DebitAmountCents:  staff.accountRUB.Balance.Mul(decimal.NewFromFloat(0.1)).Mul(staff.CurrentPrice),
			}, staff.CurrentPrice)
			if err != nil {
				return err
			}
			staff.accountUSD = resp.NewAccountFrom
			staff.accountRUB = resp.NewAccountTo
			staff.LastOpPrice = staff.CurrentPrice
			staff.IsNextOperationBuy = true
		}
	}
	return nil
}

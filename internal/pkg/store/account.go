package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/sovcomhack-inside/internal/pkg/constants"
	"github.com/sovcomhack-inside/internal/pkg/model/core"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/sovcomhack-inside/internal/pkg/store/xpgx"
)

type AccountStore interface {
	CreateAccount(ctx context.Context, account *core.Account) error
	SearchUserAccounts(ctx context.Context, opts SearchAccountsOpts) ([]core.Account, error)
	UpdateAccountBalance(ctx context.Context, accountNumber uuid.UUID, debitAmount decimal.Decimal) (acc *core.Account, txErr error)
	TransferMoney(ctx context.Context, req *dto.TransferRequestDTO) (accFrom *core.Account, accTo *core.Account, txErr error)
	SubscribeToFuntik(ctx context.Context, req *dto.SubscribeToFuntikRequest) (acc *core.Account, txErr error)
}

type SearchAccountsOpts struct {
	UserID           int64
	AccountNumbersIn []uuid.UUID
}

// CreateAccount создать счет в базе
func (s *store) CreateAccount(ctx context.Context, account *core.Account) error {
	query := builder().Insert(tableAccounts).
		Columns("number", "user_id", "currency", "for_bot").
		Values(account.Number, account.UserID, account.Currency, account.ForBot).
		Suffix("ON CONFLICT (user_id, currency, for_bot) DO NOTHING RETURNING created_at")

	if err := s.pool.QueryxRow(ctx, query).Scan(&account.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return constants.AccountAlreadyExists
		}
		return err
	}

	return nil
}

// SearchUserAccounts найти все счета по заданным параметрам
func (s *store) SearchUserAccounts(ctx context.Context, opts SearchAccountsOpts) ([]core.Account, error) {
	query := builder().Select("number", "user_id", "currency", "balance", "created_at", "for_bot").
		From(tableAccounts)
	if opts.UserID != 0 {
		query = query.Where(squirrel.Eq{"user_id": opts.UserID})
	}
	if len(opts.AccountNumbersIn) > 0 {
		query = query.Where(squirrel.Eq{"number": opts.AccountNumbersIn})
	}
	query = query.Where(squirrel.Eq{"for_bot": false})
	var accounts []core.Account

	if err := s.pool.Selectx(ctx, &accounts, query); err != nil {
		return nil, err
	}

	return accounts, nil
}

// UpdateAccountBalance изменить баланс на счете
func (s *store) UpdateAccountBalance(ctx context.Context, accountNumber uuid.UUID, deltaAmountCents decimal.Decimal) (acc *core.Account, txErr error) {
	txErr = s.withTx(ctx, func(ctx context.Context, tx Tx) error {
		var err error

		acc, err = getAccount(ctx, accountNumber, tx)
		if err != nil {
			return fmt.Errorf("getAccount err: %w", err)
		}
		if acc.Balance.Add(deltaAmountCents).LessThan(decimal.NewFromInt(0)) {
			return constants.ErrNotEnoughMoney
		}

		err = updateAccountBalance(ctx, acc, deltaAmountCents, tx)
		if err != nil {
			return fmt.Errorf("updateAccountBalance err: %w", err)
		}
		return nil
	})
	return
}

func getAccount(ctx context.Context, accountNumber uuid.UUID, executor xpgx.Executor) (*core.Account, error) {
	var account core.Account

	query := builder().Select("number", "user_id", "currency", "balance", "created_at").
		From(tableAccounts).
		Where(squirrel.Eq{"number": accountNumber})

	err := executor.Getx(ctx, &account, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account with this number not found: %s", accountNumber.String())
		}
		return nil, err
	}
	return &account, err
}

func updateAccountBalance(ctx context.Context, account *core.Account, delta decimal.Decimal, executor xpgx.Executor) error {
	query := builder().Update(tableAccounts).
		Set("balance", account.Balance.Add(delta)).
		Where(squirrel.Eq{"number": account.Number}).
		Suffix("RETURNING balance")

	if err := executor.Getx(ctx, &account.Balance, query); err != nil {
		return err
	}
	return nil
}

// TransferMoney снять деньги с одного счета, положить деньги на другой счет
func (s *store) TransferMoney(ctx context.Context, req *dto.TransferRequestDTO) (accFrom *core.Account, accTo *core.Account, txErr error) {
	txErr = s.withTx(ctx, func(ctx context.Context, tx Tx) error {
		from, err := getAccount(ctx, req.AccountFrom, tx)
		if err != nil {
			return err
		}
		if from.Balance.LessThan(req.CreditAmountCents) {
			return constants.ErrNotEnoughMoney
		}
		to, err := getAccount(ctx, req.AccountTo, tx)
		if err != nil {
			return err
		}
		err = updateAccountBalance(ctx, from, req.CreditAmountCents.Neg(), tx)
		if err != nil {
			return err
		}
		err = updateAccountBalance(ctx, to, req.DebitAmountCents, tx)
		if err != nil {
			return err
		}
		accFrom = from
		accTo = to
		return nil
	})
	return
}

func (s *store) SubscribeToFuntik(ctx context.Context, req *dto.SubscribeToFuntikRequest) (accFrom *core.Account, txErr error) {
	txErr = s.withTx(ctx, func(ctx context.Context, tx Tx) error {
		user, err := getUserByID(ctx, req.UserID, tx)
		if err != nil {
			return err
		}
		if user.SubscriptionExpiredAt != nil && (*user.SubscriptionExpiredAt).After(time.Now()) {
			return errors.New("subscription is already active")
		}
		acc, err := getAccount(ctx, req.AccountNumberFrom, tx)
		if err != nil {
			return err
		}
		if acc.Balance.LessThanOrEqual(decimal.NewFromInt(req.SubscribePriceCents)) {
			return constants.ErrNotEnoughMoney
		}
		err = updateAccountBalance(ctx, acc, decimal.NewFromInt(req.SubscribePriceCents).Neg(), tx)
		if err != nil {
			return err
		}
		err = updateSubscriptionExpiredAt(ctx, user, tx)
		if err != nil {
			return err
		}
		accFrom = acc
		return nil
	})
	return
}

func updateSubscriptionExpiredAt(ctx context.Context, user *core.User, executor xpgx.Executor) error {
	query := builder().
		Update(tableUsers).
		Set("subscription_expired_at", time.Now().AddDate(0, 1, 0)).
		Where(squirrel.Eq{"id": user.ID})
	if _, err := executor.Execx(ctx, query); err != nil {
		return err
	}
	return nil
}

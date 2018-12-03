package example

import (
	"context"

	"github.com/aneshas/tx"
)

// NewAccountService creates new account application service instance
func NewAccountService(accRepo AccountRepository) *AccountService {
	return &AccountService{
		Transactional: accRepo,
		accounts:      accRepo,
	}
}

// AccountRepository represents account repository interface
type AccountRepository interface {
	tx.Transactional // Embed interface inside of a repo or inject it separately through NewAccountService constructor

	ByID(context.Context, int64) (*Account, error)
	Save(context.Context, *Account) error
}

// AccountService represents account application service
type AccountService struct {
	tx.Transactional // A convenience so we can use .RunTx directly on a service (mimick @Transactional)

	accounts AccountRepository
}

// TransferReq DTO
type TransferReq struct {
	SrcID  int64
	DestID int64
	Amount int
}

// TransferMoney transfers the amount from source account to destination account
func (svc *AccountService) TransferMoney(ctx context.Context, req *TransferReq) error {

	// Do your authentication/acl here

	// This is the crust of the example
	// This is how we would use the tx abstraction independent of the db implementation
	return svc.RunTx(
		ctx,
		// This func will be run in a transaction
		// It is crucial for you to use the context provided by the
		// call to this function, otherwise statements will not be run in a transaction
		func(ctx context.Context) error {

			src, err := svc.accounts.ByID(ctx, req.SrcID)
			if err != nil {
				return err
			}

			dest, err := svc.accounts.ByID(ctx, req.DestID)
			if err != nil {
				return err
			}

			err = src.Withdraw(req.Amount)
			if err != nil {
				return err
			}

			dest.Deposit(req.Amount)

			err = svc.accounts.Save(ctx, src)
			if err != nil {
				return err
			}

			err = svc.accounts.Save(ctx, dest)
			if err != nil {
				return err
			}

			return nil

		},
	)
}

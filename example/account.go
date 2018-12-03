package example

import (
	"fmt"
)

// Account entity
type Account struct {
	Balance int
	UserID  int64
}

// Withdraw withdraws the amount from account if available
func (a *Account) Withdraw(amount int) error {
	if amount > a.Balance {
		return fmt.Errorf("account: not enough funds")
	}

	a.Balance -= amount

	return nil
}

// Deposit deposits amount to the account
func (a *Account) Deposit(amount int) {
	a.Balance += amount
}

package bankaccount

import (
	"fmt"
	"sync"
)

// Note that this is purely for example purposes and is not production code quality. I wrote my own
// implmentations here purely for the purpose of being able to demostrate some tests.

type Account interface {
	Balance() Money
	Deposit(Money) error
	Withdraw(Money) error
}

type SavingsAccount struct {
	balance Money
	sync.Mutex
}

type SavingsAccountOption func(*SavingsAccount)

func WithBalance(m Money) SavingsAccountOption {
	return func(s *SavingsAccount) {
		s.balance = m
	}
}

func NewSavingsAccount(opts ...SavingsAccountOption) *SavingsAccount {
	m, _ := NewMoney(USD, 0, 0)
	acct := &SavingsAccount{
		balance: m,
	}
	for _, opt := range opts {
		opt(acct)
	}
	return acct
}

func (s *SavingsAccount) Balance() Money {
	return s.balance
}

func (s *SavingsAccount) Deposit(m Money) error {
	s.Lock()
	newBalance, err := s.balance.Add(m)
	if err == nil {
		s.balance = newBalance
	}
	s.Unlock()
	return err
}

func (s *SavingsAccount) Withdraw(m Money) error {
	s.Lock()
	newBalance, err := s.balance.Subtract(m)
	if newBalance.IsNegative() {
		err = fmt.Errorf("withdrawal of %s would overdraw from balance of %s", m, s.balance)
	} else if err == nil {
		s.balance = newBalance
	}
	s.Unlock()
	return err
}

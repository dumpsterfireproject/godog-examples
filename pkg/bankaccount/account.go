package bankaccount

import (
	"fmt"
	"strconv"
	"sync"
)

// Note that this is purely for example purposes and is not production code quality. I wrote my own
// implmentations here purely for the purpose of being able to demostrate some tests.

type Account interface {
	Balance() Money
	BalanceAsCurrency(string) (Money, error)
	Deposit(Money) error
	Withdraw(Money) error
	RemittanceAddress() string
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

func (s *SavingsAccount) BalanceAsCurrency(currencyCode string) (Money, error) {
	conversion := exchangeRate{s.balance.CurrencyCode, currencyCode}
	rate, found := CurrentRates.rates[conversion]
	if !found {
		return Money{}, fmt.Errorf("currency code not found in current exchange tables")
	}
	mantissa, exponent := asExponent(rate.Units, rate.Nanos)
	m := s.balance.Multiply(int(mantissa), exponent)
	m.CurrencyCode = currencyCode

	return m, nil
}

func asExponent(units int64, nanos int32) (int, int) {
	i, _ := strconv.Atoi(fmt.Sprintf("%d%09d", units, nanos))
	return i, -9
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

func (s *SavingsAccount) RemittanceAddress() string {
	return "742 Evergreen Terrace\nSpringfield, OR"
}

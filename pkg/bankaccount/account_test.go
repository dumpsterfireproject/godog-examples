package bankaccount_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	. "github.com/dumpsterfireproject/godog-examples/pkg/bankaccount"
)

// keys
type accountKey struct{}
type errorKey struct{}

//  helper methods
func getAccount(ctx context.Context) Account {
	acct := ctx.Value(accountKey{}).(Account)
	return acct
}

// Arrange steps
func iHaveANewAccount(ctx context.Context) context.Context {
	acct := NewSavingsAccount()
	return context.WithValue(ctx, accountKey{}, acct)
}

func iHaveAnAccountWith(ctx context.Context, units int, nanos int, currency string) context.Context {
	m, _ := NewMoney(currency, int64(units), int32(nanos))
	acct := NewSavingsAccount(WithBalance(m))
	return context.WithValue(ctx, accountKey{}, acct)
}

// Act steps
func iDeposit(ctx context.Context, units int, nanos int, currency string) (context.Context, error) {
	acct := getAccount(ctx)
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return ctx, err
	}
	err = acct.Deposit(m)
	return context.WithValue(ctx, accountKey{}, acct), err
}

func iWithdraw(ctx context.Context, units int, nanos int, currency string) (context.Context, error) {
	acct := getAccount(ctx)
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return ctx, err
	}
	err = acct.Withdraw(m)
	return context.WithValue(ctx, accountKey{}, acct), err
}

func iTryToWithdraw(ctx context.Context, units int, nanos int, currency string) context.Context {
	acct := getAccount(ctx)
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return context.WithValue(ctx, errorKey{}, err)
	}
	err = acct.Withdraw(m)
	if err != nil {
		return context.WithValue(ctx, errorKey{}, err)
	}
	return context.WithValue(ctx, accountKey{}, acct)
}

// Assert steps
func theAccountBalanceIs(ctx context.Context, units int, nanos int, currency string) error {
	acct := getAccount(ctx)
	m, _ := NewMoney(currency, int64(units), int32(nanos))
	if !acct.Balance().IsEqual(m) {
		return fmt.Errorf("expected the account balance to be %s by found %s", m, acct.Balance())
	}
	return nil
}

func theTransactionShouldError(ctx context.Context) error {
	err := ctx.Value(errorKey{})
	if err == nil {
		return fmt.Errorf("the expected error was not found")
	}
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			// Add step definitions here.
			s.Step(`^I have a new account$`, iHaveANewAccount)
			s.Step(`^I have an account with (\d+)\.(\d+) ([A-Z]{3})$`, iHaveAnAccountWith)
			s.Step(`^I deposit (\d+)\.(\d+) ([A-Z]{3})$`, iDeposit)
			s.Step(`^I withdraw (\d+)\.(\d+) ([A-Z]{3})$`, iWithdraw)
			s.Step(`^I try to withdraw (\d+)\.(\d+) ([A-Z]{3})$`, iTryToWithdraw)
			s.Step(`^the account balance must be (\d+)\.(\d+) ([A-Z]{3})$`, theAccountBalanceIs)
			s.Step(`^the transaction should error$`, theTransactionShouldError)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

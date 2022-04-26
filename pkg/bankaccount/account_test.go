package bankaccount_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	. "github.com/dumpsterfireproject/godog-examples/pkg/bankaccount"
)

type AccountTestState struct {
	account   Account
	lastError error
}

func (a *AccountTestState) reset() {
	a.account = nil
	a.lastError = nil
}

// Arrange steps
func (a *AccountTestState) iHaveANewAccount() {
	a.account = NewSavingsAccount()
}

func (a *AccountTestState) iHaveAnAccountWith(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), int32(nanos))
	a.account = NewSavingsAccount(WithBalance(m))
	return err
}

// Act steps
func (a *AccountTestState) iDeposit(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return err
	}
	err = a.account.Deposit(m)
	return err
}

func (a *AccountTestState) iWithdraw(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return err
	}
	err = a.account.Withdraw(m)
	return err
}

func (a *AccountTestState) iTryToWithdraw(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), int32(nanos))
	if err != nil {
		return err
	}
	// this step does not fail if the withdrawal fails; it just stores the error for later validation
	err = a.account.Withdraw(m)
	a.lastError = err
	return nil
}

// Assert steps
func (a *AccountTestState) theAccountBalanceIs(units int, nanos int, currency string) error {
	acct := a.account
	m, _ := NewMoney(currency, int64(units), int32(nanos))
	if !acct.Balance().IsEqual(m) {
		return fmt.Errorf("expected the account balance to be %s by found %s", m, acct.Balance())
	}
	return nil
}

func (a *AccountTestState) theTransactionShouldError() error {
	if a.lastError == nil {
		return fmt.Errorf("the expected error was not found")
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	ts := &AccountTestState{}
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ts.reset() // clean the state before every scenario
		return ctx, nil
	})
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})
	// Add step definitions here.
	sc.Step(`^I have a new account$`, ts.iHaveANewAccount)
	sc.Step(`^I have an account with (\d+)\.(\d+) ([A-Z]{3})$`, ts.iHaveAnAccountWith)
	sc.Step(`^I deposit (\d+)\.(\d+) ([A-Z]{3})$`, ts.iDeposit)
	sc.Step(`^I withdraw (\d+)\.(\d+) ([A-Z]{3})$`, ts.iWithdraw)
	sc.Step(`^I try to withdraw (\d+)\.(\d+) ([A-Z]{3})$`, ts.iTryToWithdraw)
	sc.Step(`^the account balance must be (\d+)\.(\d+) ([A-Z]{3})$`, ts.theAccountBalanceIs)
	sc.Step(`^the transaction should error$`, ts.theTransactionShouldError)
}

func IntializeTestSuite(sc *godog.TestSuiteContext) {
	sc.BeforeSuite(func() {
		// do any set up here one time before the entire test suite runs, e.g., create a database connection pool
		// that would be too expensive to do before each scenario.
	})
	sc.AfterSuite(func() {
		// do any clean up after the entire test suite is done executiong.
	})
	sc.ScenarioContext().StepContext().Before(func(ctx context.Context, st *godog.Step) (context.Context, error) {
		return ctx, nil
	})
	sc.ScenarioContext().StepContext().After(func(ctx context.Context, st *godog.Step, status godog.StepResultStatus, err error) (context.Context, error) {
		return ctx, nil
	})
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: IntializeTestSuite,
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

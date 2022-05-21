package bankaccount_test

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
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
	m, err := NewMoney(currency, int64(units), convertToNanos(nanos))
	a.account = NewSavingsAccount(WithBalance(m))
	return err
}

// Act steps
func (a *AccountTestState) iDeposit(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), convertToNanos(nanos))
	if err != nil {
		return err
	}
	err = a.account.Deposit(m)
	return err
}

func (a *AccountTestState) iWithdraw(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), convertToNanos(nanos))
	if err != nil {
		return err
	}
	err = a.account.Withdraw(m)
	return err
}

func (a *AccountTestState) iTryToWithdraw(units int, nanos int, currency string) error {
	m, err := NewMoney(currency, int64(units), convertToNanos(nanos))
	if err != nil {
		return err
	}
	// this step does not fail if the withdrawal fails; it just stores the error for later validation
	err = a.account.Withdraw(m)
	a.lastError = err
	return nil
}

type transaction struct {
	isWithdrawal bool
	money        Money
}

func (a *AccountTestState) iProcessTheFollowingTransations(table *godog.Table) error {
	transactions := []transaction{}
	// first row is header row, so skip it
	for n, row := range table.Rows {
		if n > 0 {
			if len(row.Cells) < 2 {
				return fmt.Errorf("too few columns")
			}
			units, err := strconv.Atoi(row.Cells[1].Value)
			if err != nil {
				return err
			}
			money, err := NewMoney(USD, int64(units), 0)
			if err != nil {
				return err
			}
			t := transaction{strings.ToLower(row.Cells[0].Value) == "withdrawal", money}
			transactions = append(transactions, t)
		}
	}
	wg := sync.WaitGroup{}
	wg.Add(len(transactions))
	for _, t := range transactions {
		tr := t
		go func() {
			defer wg.Done()
			if tr.isWithdrawal {
				a.account.Withdraw(tr.money)
			} else {
				a.account.Deposit(tr.money)
			}
		}()
	}
	wg.Wait()
	return nil
}

// Assert steps
func (a *AccountTestState) theAccountBalanceIs(units int, nanos int, currency string) error {
	acct := a.account
	m, _ := NewMoney(currency, int64(units), convertToNanos(nanos))
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

func (a *AccountTestState) theAccountBalanceMustConvertToUSD(input string) error {
	tokens := strings.Split(input, ".")
	if len(tokens) > 2 {
		return fmt.Errorf("invalid format for dollars")
	}
	units, err := strconv.Atoi(tokens[0])
	if err != nil {
		return err
	}
	nanos := 0
	if len(tokens) > 1 {
		nanos, err = strconv.Atoi(tokens[1])
		if err != nil {
			return err
		}
	}

	expectedDollars, _ := NewMoney(USD, int64(units), convertToNanos(nanos))
	actualDollars, err := a.account.BalanceAsCurrency(USD)
	if err != nil {
		return err
	}
	if !actualDollars.IsEqual(expectedDollars) {
		return fmt.Errorf("expected the account balance to be %s but found %s", expectedDollars, actualDollars)
	}
	return nil
}

func (a *AccountTestState) theRemittanceAddressMustBe(input *godog.DocString) error {
	if a.account.RemittanceAddress() != input.Content {
		return fmt.Errorf("expected %s but found %s", input.Content, a.account.RemittanceAddress())
	}
	// return nil
	return fmt.Errorf("invalid remittance")
}

// helper functions

// if the step has something like 1.25, the 25 is really 250000000 nanos
// this function handles that
func convertToNanos(i int) int32 {
	multiplier := 100000000
	remainingDigits := i
	for {
		remainingDigits = remainingDigits / 10
		if remainingDigits == 0 {
			break
		}
		multiplier = multiplier / 10
	}
	return int32(i * multiplier)
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
	sc.Step(`^I process the following transations:$`, ts.iProcessTheFollowingTransations)
	sc.Step(`^the account balance must be (\d+)\.(\d+) ([A-Z]{3})$`, ts.theAccountBalanceIs)
	sc.Step(`^the transaction should error$`, ts.theTransactionShouldError)
	sc.Step(`^the account balance must convert to (\d+) USD$`, ts.theAccountBalanceMustConvertToUSD)
	sc.Step(`^the remittance address must be$`, ts.theRemittanceAddressMustBe)
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

// os.Getenv("GODOG_TAGS")
var tags = flag.String("godog.tags", "", "tags to execute")
var format = flag.String("godog.format", "pretty", "format")

var opts = &godog.Options{
	Paths: []string{"features"},
}

func TestFeatures(t *testing.T) {
	opts.TestingT = t
	opts.Tags = *tags
	opts.Format = *format

	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: IntializeTestSuite,
		Options:              opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

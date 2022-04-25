package bankaccount

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

type additionTestCase struct {
	input1        Money
	input2        Money
	result        Money
	expectedError error
}

func newTestCase(input1 Money, input2 Money, result Money) additionTestCase {
	return additionTestCase{input1: input1, input2: input2, result: result, expectedError: nil}
}

func (tc additionTestCase) Name(i int) string {
	return fmt.Sprintf("Example %d: %s%d.%d, %s%d.%d", i,
		tc.input1.CurrencyCode, tc.input1.Units, tc.input1.Nanos,
		tc.input2.CurrencyCode, tc.input2.Units, tc.input2.Nanos)
}

func (tc additionTestCase) RunAdd(t *testing.T) {
	is := is.NewRelaxed(t)
	actual, err := tc.input1.Add(tc.input2)
	if tc.expectedError != nil {
		if err == nil {
			t.Errorf("did not get the expected error")
		} else {
			is.Equal(err.Error(), tc.expectedError.Error())
		}
	} else {
		is.NoErr(err)
	}
	is.Equal(actual.CurrencyCode, tc.result.CurrencyCode) // CurrencyCode
	is.Equal(actual.Units, tc.result.Units)               // Units
	is.Equal(actual.Nanos, tc.result.Nanos)               // Nanos
}

func (tc additionTestCase) RunSubtract(t *testing.T) {
	is := is.NewRelaxed(t)
	actual, err := tc.input1.Subtract(tc.input2)
	if tc.expectedError != nil {
		if err == nil {
			t.Errorf("did not get the expected error")
		} else {
			is.Equal(err.Error(), tc.expectedError.Error())
		}
	} else {
		is.NoErr(err)
	}
	is.Equal(actual.CurrencyCode, tc.result.CurrencyCode) // CurrencyCode
	is.Equal(actual.Units, tc.result.Units)               // Units
	is.Equal(actual.Nanos, tc.result.Nanos)               // Nanos
}

func TestAdd(t *testing.T) {
	testCases := []additionTestCase{
		{Money{USD, 0, 0}, Money{CAD, 0, 0}, Money{}, fmt.Errorf("you must convert values to common currency code using current exchange rates before adding")},
		newTestCase(Money{USD, 0, 0}, Money{USD, 0, 0}, Money{USD, 0, 0}),
		newTestCase(Money{USD, 0, 0}, Money{USD, 1, 1}, Money{USD, 1, 1}),
		newTestCase(Money{USD, 0, 0}, Money{USD, -1, -1}, Money{USD, -1, -1}),
		newTestCase(Money{USD, 1, 1}, Money{USD, -1, -1}, Money{USD, 0, 0}),
		newTestCase(Money{USD, -1, -1}, Money{USD, 1, 1}, Money{USD, 0, 0}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, 1, 500000000}, Money{USD, 3, 0}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, 1, 600000000}, Money{USD, 3, 100000000}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, -1, -500000000}, Money{USD, -3, 0}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, -1, -600000000}, Money{USD, -3, -100000000}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, 1, 500000000}, Money{USD, 0, 0}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, 1, 600000000}, Money{USD, 0, 100000000}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, -1, -500000000}, Money{USD, 0, 0}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, -1, -600000000}, Money{USD, 0, -100000000}),
		newTestCase(Money{USD, 1, 1}, Money{USD, -10, -999999999}, Money{USD, -9, -999999998}),
	}
	for i, tc := range testCases {
		t.Run(tc.Name(i), tc.RunAdd)
	}
}

func TestSubtract(t *testing.T) {
	testCases := []additionTestCase{
		{Money{USD, 0, 0}, Money{CAD, 0, 0}, Money{}, fmt.Errorf("you must convert values to common currency code using current exchange rates before adding")},
		newTestCase(Money{USD, 0, 0}, Money{USD, 0, 0}, Money{USD, 0, 0}),
		newTestCase(Money{USD, 0, 0}, Money{USD, 1, 1}, Money{USD, -1, -1}),
		newTestCase(Money{USD, 0, 0}, Money{USD, -1, -1}, Money{USD, 1, 1}),
		newTestCase(Money{USD, 1, 1}, Money{USD, -1, -1}, Money{USD, 2, 2}),
		newTestCase(Money{USD, -1, -1}, Money{USD, 1, 1}, Money{USD, -2, -2}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, 1, 500000000}, Money{USD, 0, 0}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, 1, 600000000}, Money{USD, 0, -100000000}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, -1, -500000000}, Money{USD, 0, 0}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, -1, -600000000}, Money{USD, 0, 100000000}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, 1, 500000000}, Money{USD, -3, 0}),
		newTestCase(Money{USD, -1, -500000000}, Money{USD, 1, 600000000}, Money{USD, -3, -100000000}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, -1, -500000000}, Money{USD, 3, 0}),
		newTestCase(Money{USD, 1, 500000000}, Money{USD, -1, -600000000}, Money{USD, 3, 100000000}),
		newTestCase(Money{USD, 1, 1}, Money{USD, 10, 999999999}, Money{USD, -9, -999999998}),
	}
	for i, tc := range testCases {
		t.Run(tc.Name(i), tc.RunSubtract)
	}
}

type multiplicationTestCase struct {
	money    Money
	mantissa int
	exponent int
	expected Money
}

func (tc multiplicationTestCase) Name(i int) string {
	return fmt.Sprintf("Example %d: %s%d.%d * %dx10^%d",
		i, tc.money.CurrencyCode, tc.money.Units, tc.money.Nanos, tc.mantissa, tc.exponent)
}

func (tc multiplicationTestCase) Run(t *testing.T) {
	is := is.NewRelaxed(t)
	actual := tc.money.Multiply(tc.mantissa, tc.exponent)
	is.Equal(actual.CurrencyCode, tc.expected.CurrencyCode) // CurrencyCode
	is.Equal(actual.Units, tc.expected.Units)               // Units
	is.Equal(actual.Nanos, tc.expected.Nanos)               // Nanos
}

func TestMultiply(t *testing.T) {
	testCases := []multiplicationTestCase{
		// multiplication by 0
		{Money{USD, 0, 0}, 0, 0, Money{USD, 0, 0}},
		{Money{USD, 1, 1}, 0, 0, Money{USD, 0, 0}},
		{Money{USD, -1, -1}, 0, 0, Money{USD, 0, 0}},
		// multiplication by 1
		{Money{USD, 0, 0}, 1, 0, Money{USD, 0, 0}},
		{Money{USD, 1, 1}, 1, 0, Money{USD, 1, 1}},
		{Money{USD, -1, -1}, 1, 0, Money{USD, -1, -1}},
		// multiplication by -1
		{Money{USD, 0, 0}, -1, 0, Money{USD, 0, 0}},
		{Money{USD, 1, 1}, -1, 0, Money{USD, -1, -1}},
		{Money{USD, -1, -1}, -1, 0, Money{USD, 1, 1}},
		// multiplication by whole number
		{Money{USD, 2, 300000000}, 2, 0, Money{USD, 4, 600000000}},
		{Money{USD, 2, 600000000}, 2, 0, Money{USD, 5, 200000000}},
		{Money{USD, -2, -300000000}, 2, 0, Money{USD, -4, -600000000}},
		{Money{USD, -2, -600000000}, 2, 0, Money{USD, -5, -200000000}},
		{Money{USD, 2, 300000000}, -2, 0, Money{USD, -4, -600000000}},
		{Money{USD, 2, 600000000}, -2, 0, Money{USD, -5, -200000000}},
		// multiplication by fraction
		{Money{USD, 4, 600000000}, 5, -1, Money{USD, 2, 300000000}},
		{Money{USD, 5, 200000000}, 5, -1, Money{USD, 2, 600000000}},
		{Money{USD, -4, -600000000}, 5, -1, Money{USD, -2, -300000000}},
		{Money{USD, -5, -200000000}, 5, -1, Money{USD, -2, -600000000}},
		{Money{USD, -4, -600000000}, -5, -1, Money{USD, 2, 300000000}},
		{Money{USD, -5, -200000000}, -5, -1, Money{USD, 2, 600000000}},
		// multiplication by > 1.00, < -1.00
		{Money{USD, 4, 600000000}, 15, -1, Money{USD, 6, 900000000}},
		{Money{USD, 5, 200000000}, 15, -1, Money{USD, 7, 800000000}},
		{Money{USD, -4, -600000000}, 15, -1, Money{USD, -6, -900000000}},
		{Money{USD, -5, -200000000}, 15, -1, Money{USD, -7, -800000000}},
		{Money{USD, -4, -600000000}, -15, -1, Money{USD, 6, 900000000}},
		{Money{USD, -5, -200000000}, -15, -1, Money{USD, 7, 800000000}},
		// rounding fractional nanos
		{Money{USD, 0, 1}, 6, -1, Money{USD, 0, 1}},
		{Money{USD, 0, 1}, 5, -1, Money{USD, 0, 1}},
		{Money{USD, 0, 1}, 2, -1, Money{USD, 0, 0}},
	}
	for i, tc := range testCases {
		t.Run(tc.Name(i), tc.Run)
	}
}

func TestIsNegative(t *testing.T) {
	testCases := []struct {
		money    Money
		expected bool
	}{
		{Money{USD, 0, 1}, false},
		{Money{USD, 1, 1}, false},
		{Money{USD, 1, 0}, false},
		{Money{USD, 0, 0}, false},
		{Money{USD, 0, -1}, true},
		{Money{USD, -1, -1}, true},
		{Money{USD, -1, 0}, true},
	}
	for _, tc := range testCases {
		is := is.New(t)
		is.Equal(tc.money.IsNegative(), tc.expected)
	}
}

package bankaccount

const (
	USD = "USD"
	CAD = "CAD"
	CNY = "CNY"
	EUR = "EUR"
)

type exchangeRate struct {
	from string
	to   string
}

type ExchangeRates struct {
	rates map[exchangeRate]Money
}

// func (e ExchangeRates) ToUSD(money Money) (Money, error) {
// 	conversion := exchangeRate{from: money.CurrencyCode, to: USD}
// 	currentRate, exists := e.rates[conversion]
// 	if !exists {
// 		return Money{}, fmt.Errorf("current exchange rate not found for %s", money.CurrencyCode)
// 	}
// 	m, err := money.Multiply(currentRate.Units, currentRate.Nanos)
// 	return m, err
// }

var CurrentRates = ExchangeRates{
	rates: map[exchangeRate]Money{
		{CAD, USD}: {USD, 0, 800000000},
		{CNY, USD}: {USD, 0, 160000000},
		{EUR, USD}: {USD, 1, 80000000},
	},
}

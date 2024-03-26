package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	JPN = "JPN"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, JPN:
		return true
	}
	return false
}

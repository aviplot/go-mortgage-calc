package mortgage

import "github.com/aviplot/go-finance-math/financial"
import "github.com/shopspring/decimal"

// Flow is full cashflow of mortgage
type Flow struct {
	Index           int64
	Date            financial.Date
	Principal       decimal.Decimal
	Interest        decimal.Decimal
	Payment         decimal.Decimal
	PrincipalLeft   decimal.Decimal
	KnownIndex      decimal.Decimal
	PaymentFactored decimal.Decimal
}

func (f Flow) ToCashFlow() (cf financial.CashFlow) {
	cf.Flow = f.Payment.InexactFloat64()
	cf.Date = f.Date
	return
}

func (f Flow) ToCashFlowIndexed() (cf financial.CashFlow) {
	cf.Flow = f.PaymentFactored.InexactFloat64()
	cf.Date = f.Date
	return
}

package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
)

// indexZero return always zero.
type indexZero struct {
	baseIndex     decimal.Decimal
	baseIndexDate financial.Date
}

func (i indexZero) GetIndex(d financial.Date) decimal.Decimal {
	return decimal.Zero
}

func (i indexZero) GetBaseIndex() (f financial.Date, d decimal.Decimal) {
	return
}

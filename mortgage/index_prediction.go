package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
)

// IndexPrediction allow prediction of future index
type IndexPrediction interface {
	GetIndex(date financial.Date) decimal.Decimal
	GetBaseIndex() (financial.Date, decimal.Decimal)
}

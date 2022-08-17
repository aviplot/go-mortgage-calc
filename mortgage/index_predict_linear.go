package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
)

// indexGrowthLinear return linear growth of the index.
type indexGrowthLinear struct {
	baseIndex     decimal.Decimal
	baseIndexDate financial.Date
	curve         decimal.Decimal
}

// NewLinearGrowthIP create new index predictor (linear curve)
func NewLinearGrowthIP(bd financial.Date, bi, curve decimal.Decimal) (i indexGrowthLinear) {
	i.baseIndex = bi
	i.baseIndexDate = bd
	i.curve = curve
	return
}

func (i indexGrowthLinear) GetIndex(d financial.Date) decimal.Decimal {
	dd := decimal.NewFromInt(d.DaysFrom(i.baseIndexDate))
	return dd.Mul(i.curve).Add(i.baseIndex).RoundBank(1)
}

func (i indexGrowthLinear) GetBaseIndex() (f financial.Date, d decimal.Decimal) {
	return i.baseIndexDate, i.baseIndex
}

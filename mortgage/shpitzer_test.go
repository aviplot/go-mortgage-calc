package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
	"testing"
)

func TestShpitzer(t *testing.T) {

	//t CalculationType, a, r decimal.Decimal, m int64, ad, pd financial.Date
	a := decimal.NewFromFloat(20000.0)
	r := decimal.NewFromFloat(0.0125)
	m := 18
	ad := financial.NewDateFromFormattedString("2000-01-15")
	pd := financial.NewDateFromFormattedString("2000-02-10")
	ip := NewLinearGrowthIP(financial.NewDateFromFormattedString("1999-01-15"), decimal.NewFromFloat(106.2), decimal.NewFromFloat(0.005))
	cf, e := CreateFlowTable(Shpitzer, a, r, int64(m), ad, pd, ip)
	if e != nil {
		t.Fatalf("Error: %v", e)
	}
	t.Logf("Cash flow: %v", cf)
	irr, _ := cf.Irr()
	t.Logf("IRR: %v", irr)
}

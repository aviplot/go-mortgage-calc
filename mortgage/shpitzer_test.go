package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
	"log"
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
	kf := make(financial.CashFlowTab, 0)
	kf = append(kf, financial.CashFlow{
		Date: financial.NewDateFromFormattedString("2005-02-10"),
		Flow: 1500.20,
	})
	cf, e := CreateFlowTable(Shpitzer, a, r, int64(m), ad, pd, ip, kf)
	if e != nil {
		t.Fatalf("Error: %v", e)
	}
	t.Logf("Cash flow: %v", cf)
	expected := decimal.NewFromFloat(0.2099077)
	irr, _ := cf.Irr()
	if !irr.RoundDown(8).Equals(expected) {
		log.Fatalf("Got: %v, but expected: %v", irr.RoundDown(8), expected)
	}

}

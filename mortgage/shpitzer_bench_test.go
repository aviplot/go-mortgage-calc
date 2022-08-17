package mortgage

import (
	"fmt"
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
	"testing"
)

func BenchmarkShpitzerCashFlow(b *testing.B) {

	//t CalculationType, a, r decimal.Decimal, m int64, ad, pd financial.Date
	a := decimal.NewFromFloat(20000.0)
	r := decimal.NewFromFloat(0.0125)
	//m := 360
	ad := financial.NewDateFromFormattedString("2000-01-15")
	pd := financial.NewDateFromFormattedString("2000-02-10")
	ip := NewLinearGrowthIP(financial.NewDateFromFormattedString("1999-01-15"), decimal.NewFromFloat(106.2), decimal.NewFromFloat(0.005))
	kf := make(financial.CashFlowTab, 0)
	kf = append(kf, financial.CashFlow{
		Date: financial.NewDateFromFormattedString("2005-02-10"),
		Flow: 1500.20,
	})

	b.ResetTimer()

	for m := 12; m <= 360; m += 6 { // From one year, up to 30 years, calculation for each half month
		b.Run(fmt.Sprintf("Months_%d", m), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cf, e := CreateFlowTable(Shpitzer, a, r, int64(m), ad, pd, ip, kf)
				if e != nil {
					b.Fatalf("Error: %v", e)
				}

				// expected := decimal.NewFromFloat(0.2099077) // 18 months
				// expected := decimal.NewFromFloat(0.16746812) // 360 months
				irr, _ := cf.Irr()
				//if !irr.RoundDown(8).Equals(expected) {
				//	b.Fatalf("Got: %v, but expected: %v", irr.RoundDown(8), expected)
				//}
				if irr.IsNegative() {
					b.Fatalf("Irr is less the 0, (%v)", irr)
				}
			}
		})
	}

}

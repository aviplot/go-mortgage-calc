package main

import (
	"fmt"
	"github.com/aviplot/go-finance-math/financial"
	"github.com/aviplot/go-mortgage-calc/mortgage"
	"github.com/shopspring/decimal"
	"time"
	"os"
	"github.com/aviplot/go-codegen-datastruct/datastruct"
	//"sync"
)

const (
	amountFrom = 20000
	amountTo   = 500000
	amountGap  = 5000

	monthlyRate       = 0.0125

	monthsFrom = 18
	monthsTo   = 360
	monthsGap  = 6
)

func saveStringToFile(data, path string) error {
    return os.WriteFile(path, []byte(data), 0644)
}

func main() {
    data := make(map[string]decimal.Decimal, 0) // index is: amount,pmt -> value is irr


	r := decimal.NewFromFloat(monthlyRate)
	ad := financial.NewDateFromFormattedString("2000-01-15")
	pd := financial.NewDateFromFormattedString("2000-02-10")
	ip := mortgage.NewLinearGrowthIP(financial.NewDateFromFormattedString("1999-01-15"), decimal.NewFromFloat(106.2), decimal.NewFromFloat(0.005))
	kf := make(financial.CashFlowTab, 0)
	kf = append(kf, financial.CashFlow{
		Date: financial.NewDateFromFormattedString("2005-02-10"),
		Flow: 1500.20,
	})
	calcCounter := 0

    //var wg sync.WaitGroup
    var key string
	start := time.Now()
	for amount := amountFrom; amount < amountTo; amount += amountGap {
	    //wg.Add(1)

        //go func() {
            //defer wg.Done()

            a := decimal.NewFromInt(int64(amount))
            for m:= monthsFrom; m<monthsTo; m+=monthsGap {
                cft, e := mortgage.CreateFlowTable(mortgage.Shpitzer, a, r, int64(m), ad, pd, ip, kf)
                if e != nil {
                    fmt.Println("Error....")
                }
                // Save to variable
                irr, _ := cft.Irr()
                if irr.GreaterThanOrEqual(decimal.NewFromFloat(0.5)) {
                    fmt.Println(irr)
                }
                key = fmt.Sprintf("%v-%v", amount, m)
                data[key] = irr
                calcCounter++
            }
       // }()

	}
	//wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Execution takes", elapsed)
	fmt.Println("Calculation: ", calcCounter)
	fmt.Println("Saving to file...")

    ds, _ := datastruct.NewDataStruct(MapDataType, d, "DataVariableName")
	saveStringToFile(ds.ExportToCode(), "data.go.txt")
}

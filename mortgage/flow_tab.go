package mortgage

import (
	"errors"
	"fmt"
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
	"sort"
	"strings"
)

// FlowTab is table of flows, represent full cash flow
type FlowTab []Flow

// Common errors.
var (
	ErrParametersMissing = errors.New("invalid data, calculation failed due to missing parameters")
	ErrCalculationFailed = errors.New("calculation failed")
)

// CreateShpitzerFlowTable calculate and create cashflow.
// a = amount
// r = monthly rate (note: not yearly rate, divide by 12 if needed)
// m = number of months
// ad = receiving the loan -> date
// pd = start paying -> date
// ed = end paying -> date
// bi = Base index
// kf = Known flows (down payment, fees and so on...)
func CreateFlowTable(t CalculationType, a, r decimal.Decimal, m int64, ad, pd financial.Date, ip IndexPrediction, kf financial.CashFlowTab) (ft FlowTab, e error) {
	var flowCreation CashFlowCreator
	if t == Shpitzer {
		flowCreation = GetShpitzerCashflowInstance(a, r, m, ad, pd, ip, kf)
	}

	if flowCreation == nil {
		e = ErrParametersMissing
		return
	}

	return flowCreation.NewCashFlowTable()
}

// String output is the data as CSV
func (ft FlowTab) String() string {
	all := make([]string, len(ft)+1)

	// First line is titles.
	all[0] = "Index, Date, Principal, Interest, Payment, PrincipalLeft, KnownIndex, PaymentFactored"

	for i, f := range ft {
		all[i+1] = fmt.Sprintf("%v, %v, %v, %v, %v, %v, %v, %v", f.Index, f.Date, f.Principal, f.Interest, f.Payment, f.PrincipalLeft, f.KnownIndex, f.PaymentFactored)
	}
	return strings.Join(all, "\n")
}

func (ft FlowTab) ToCashFlowTab() (cft financial.CashFlowTab) {
	for _, f := range ft {
		cft = append(cft, f.ToCashFlow())
	}
	return
}

func (ft FlowTab) ToCashFlowTabIndexed() (cft financial.CashFlowTab) {
	for _, f := range ft {
		cft = append(cft, f.ToCashFlowIndexed())
	}
	return
}

func (ft FlowTab) Irr() (decimal.Decimal, error) {
	irr, err := financial.Xirr(ft.ToCashFlowTab())
	return decimal.NewFromFloat(irr), err
}

func (ft FlowTab) IrrIndexed() (decimal.Decimal, error) {
	irr, err := financial.Xirr(ft.ToCashFlowTabIndexed())
	return decimal.NewFromFloat(irr), err
}

func (ft FlowTab) Rearrange() (result FlowTab) {
	result = ft.OrderByDate()
	for i := range result {
		result[i].Index = int64(i)
	}
	return
}

// Len impl "Interface" to support sorting, using sort.Sort.
func (ft FlowTab) Len() int {
	return len(ft)
}

// Swap impl "Interface" to support sorting, using sort.Sort.
func (ft FlowTab) Swap(i, j int) {
	ft[i], ft[j] = ft[j], ft[i]
}

// Less impl "Interface" to support sorting, using sort.Sort.
func (ft FlowTab) Less(i, j int) bool {
	return ft[i].Date.Date.Before(ft[j].Date.Date)
}

func (ft FlowTab) OrderByDate() (r FlowTab) {
	r = ft
	sort.Sort(r)
	return
}

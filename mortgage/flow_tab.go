package mortgage

import (
	"errors"
	"fmt"
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
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
func CreateFlowTable(t CalculationType, a, r decimal.Decimal, m int64, ad, pd financial.Date, ip IndexPrediction) (ft FlowTab, e error) {
	var flowCreation CashFlowCreator
	if t == Shpitzer {
		flowCreation = GetShpitzerCashflowInstance(a, r, m, ad, pd, ip)
	}

	if flowCreation == nil {
		e = ErrParametersMissing
		return
	}

	return flowCreation.NewCashFlowTable()
}

func (ft FlowTab) String() string {
	l := len(ft)
	all := make([]string, l)
	for i, f := range ft {
		all[i] = fmt.Sprintf("Index: %v, Date: %v, Principal: %v, Interest: %v, Payment: %v, PrincipalLeft: %v, KnownIndex: %v, PaymentFactored: %v", f.Index, f.Date, f.Principal, f.Interest, f.Payment, f.PrincipalLeft, f.KnownIndex, f.PaymentFactored)
	}
	return strings.Join(all, "\n")
}

func (ft FlowTab) ToCashFlowTab() (cft financial.CashFlowTab) {
	for _, f := range ft {
		cft = append(cft, f.ToCashFlow())
	}
	return
}

func (ft FlowTab) Irr() (decimal.Decimal, error) {
	irr, err := financial.Xirr(ft.ToCashFlowTab())
	return decimal.NewFromFloat(irr), err
}

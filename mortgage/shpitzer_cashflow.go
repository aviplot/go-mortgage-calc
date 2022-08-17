package mortgage

import (
	"github.com/aviplot/go-finance-math/financial"
	"github.com/shopspring/decimal"
)

type ShpitzerCashflow struct {
	Amount          decimal.Decimal
	Months          int64
	MonthlyRate     decimal.Decimal
	CashRecieveDate financial.Date
	StartPayingDate financial.Date
	IndexPredict    IndexPrediction
	KnownCashFlows  financial.CashFlowTab
}

// GetShpitzerCashflowInstance Creates instance.
// Parameters are:
// a : Amount of money
// r : Rate (monthly rate, not yearly)
// m : How many months
// cashDate : Date the money is taken from the bank
// fromDate : The first payment (to the bank) date
// ip : instance of index prediction. can be null if not relevant.
// kf : Known cash flows. Optional: additional known flows, like: down payment, taxes and so on.
func GetShpitzerCashflowInstance(a, r decimal.Decimal, m int64, cashDate, fromDate financial.Date, ip IndexPrediction, kf financial.CashFlowTab) (result *ShpitzerCashflow) {
	result = new(ShpitzerCashflow)
	result.Amount = a
	result.Months = m
	result.MonthlyRate = r
	result.CashRecieveDate = cashDate
	result.StartPayingDate = fromDate
	result.IndexPredict = ip
	result.KnownCashFlows = kf
	return
}

func (s ShpitzerCashflow) getFirstFlow() (fl Flow) {
	fl.Index = 0
	fl.Date = s.CashRecieveDate
	fl.Principal = decimal.Zero
	fl.Interest = decimal.Zero
	fl.Payment = s.Amount.Neg()
	fl.KnownIndex = decimal.Zero
	fl.PaymentFactored = decimal.Zero
	return
}

func (s ShpitzerCashflow) getKnownFlows() (flt FlowTab) {
	var fl Flow
	var currentIndex decimal.Decimal
	flt = make(FlowTab, len(s.KnownCashFlows))
	_, baseIndex := s.IndexPredict.GetBaseIndex()

	for i, kf := range s.KnownCashFlows {
		fl.Index = int64(i)
		fl.Date = kf.Date
		fl.Payment = decimal.NewFromFloat(kf.Flow).RoundBank(2)

		if !baseIndex.IsZero() && s.IndexPredict != nil {
			currentIndex = s.IndexPredict.GetIndex(kf.Date)
			fl.PaymentFactored = currentIndex.Div(baseIndex).Mul(fl.Payment).RoundBank(2)
			fl.KnownIndex = currentIndex.RoundBank(2)
		}
		flt[i] = fl
	}
	return
}

func (s ShpitzerCashflow) getPaymentFlows() (flt FlowTab) {

	var (
		fromInterest  decimal.Decimal
		fromPrincipal decimal.Decimal
		currentIndex  decimal.Decimal
		fl            Flow
	)
	flt = make(FlowTab, s.Months)

	// Creating the cash flow.
	pmt := financial.Pmt(s.MonthlyRate.InexactFloat64(), s.Months, s.Amount.Neg().InexactFloat64(), 0, false)
	decPmt := decimal.NewFromFloat(pmt)

	paymDate := s.StartPayingDate
	AmountLeft := s.Amount
	_, baseIndex := s.IndexPredict.GetBaseIndex()
	for i := 0; int64(i) < s.Months; i++ {
		//indexDec := decimal.NewFromInt(int64(i))
		fromInterest = AmountLeft.Mul(s.MonthlyRate)
		fromPrincipal = decPmt.Sub(fromInterest)

		fl.Index = int64(i)
		fl.Date = paymDate
		fl.Payment = decPmt.RoundBank(2)
		fl.Interest = fromInterest.RoundBank(2)
		fl.Principal = fromPrincipal.RoundBank(2)
		fl.PrincipalLeft = AmountLeft.Sub(fromPrincipal).RoundBank(2)

		if !baseIndex.IsZero() && s.IndexPredict != nil {
			currentIndex = s.IndexPredict.GetIndex(paymDate)
			fl.PaymentFactored = currentIndex.Div(baseIndex).Mul(fl.Payment).RoundBank(2)
			fl.KnownIndex = currentIndex.RoundBank(2)
		}

		flt[i] = fl

		paymDate = paymDate.AddMonth()
		AmountLeft = AmountLeft.Sub(fromPrincipal)
	}
	return
}

func (s ShpitzerCashflow) NewCashFlowTable() (ft FlowTab, e error) {
	// Validate input
	if s.Amount.LessThanOrEqual(decimal.Zero) || s.Months < 0 {
		e = ErrParametersMissing
		return
	}

	ft = make(FlowTab, 0)

	ft = append(ft, s.getFirstFlow())
	ft = append(ft, s.getPaymentFlows()...)
	ft = append(ft, s.getKnownFlows()...)

	ft = ft.Rearrange()

	return
}

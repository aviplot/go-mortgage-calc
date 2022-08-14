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
}

func GetShpitzerCashflowInstance(a, r decimal.Decimal, m int64, cashDate, fromDate financial.Date, ip IndexPrediction) (result *ShpitzerCashflow) {
	result = new(ShpitzerCashflow)
	result.Amount = a
	result.Months = m
	result.MonthlyRate = r
	result.CashRecieveDate = cashDate
	result.StartPayingDate = fromDate
	result.IndexPredict = ip
	return
}

func (s ShpitzerCashflow) NewCashFlowTable() (ft FlowTab, e error) {
	// Validate input
	if s.Amount.LessThanOrEqual(decimal.Zero) || s.Months < 0 {
		e = ErrParametersMissing
		return
	}

	// Creating the cash flow.
	pmt := financial.Pmt(s.MonthlyRate.InexactFloat64(), s.Months, s.Amount.Neg().InexactFloat64(), 0, false)
	decPmt := decimal.NewFromFloat(pmt)

	ft = make(FlowTab, s.Months+1)
	var fl Flow

	// First flow is the cash payment
	fl.Index = 0
	fl.Date = s.CashRecieveDate
	fl.Principal = decimal.Zero
	fl.Interest = decimal.Zero
	fl.Payment = s.Amount.Neg()
	fl.KnownIndex = decimal.Zero
	fl.PaymentFactored = decimal.Zero
	ft[0] = fl

	var (
		fromInterest decimal.Decimal
		fromPrincipal decimal.Decimal
		currentIndex decimal.Decimal
	)

	paymDate := s.StartPayingDate
	AmountLeft := s.Amount
	_, baseIndex := s.IndexPredict.GetBaseIndex()
	for i := 1; int64(i) <= s.Months; i++ {
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

		ft[i] = fl

		paymDate = paymDate.AddMonth()
		AmountLeft = AmountLeft.Sub(fromPrincipal)
	}

	return
}

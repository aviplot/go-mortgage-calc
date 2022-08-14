package mortgage

type CashFlowCreator interface {
	NewCashFlowTable() (FlowTab, error)
}

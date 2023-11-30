package common

type Message struct {
	PaymentId      string           `json:"paymentId"`
	TenantId       int              `json:"tenantId"`
	BankId         int              `json:"bankId"`
	Tier           string           `json:"tier"`
	BaseCurrency   string           `json:"baseCurrency"`
	TargetCurrency string           `json:"targetCurrency"`
	StartedOn      int64            `json:"startedOn"`
	CompletedOn    int64            `json:"completedOn"`
	Stages         []ProcessingInfo `json:"stages"`
}

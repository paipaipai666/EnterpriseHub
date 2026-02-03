package domain

type PaymentStatus string

const (
	PaymentInit       PaymentStatus = "INIT"
	PaymentSuccess    PaymentStatus = "SUCCESS"
	PaymentFailed     PaymentStatus = "FAILED"
	PaymentProcessing PaymentStatus = "PROCESSING"
)

func (ps *PaymentStatus) ToStatus() int {
	if ps == nil {
		return 0
	}
	switch *ps {
	case PaymentInit:
		return 1
	case PaymentSuccess:
		return 2
	case PaymentFailed:
		return 3
	case PaymentProcessing:
		return 4
	default:
		return 0
	}
}

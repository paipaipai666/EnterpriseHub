package domain

type PaymentChannel string

const (
	ALIPAY  PaymentChannel = "ALIPAY"
	WECHAT  PaymentChannel = "WECHAT"
	BALANCE PaymentChannel = "BALANCE"
)

func TakeChannel(channel string) PaymentChannel {
	switch channel {
	case "PAYMENT_METHOD_ALIPAY":
		return ALIPAY
	case "PAYMENT_METHOD_WECHAT":
		return WECHAT
	case "PAYMENT_METHOD_BALANCE":
		return BALANCE
	default:
		return BALANCE
	}
}

func (pc *PaymentChannel) ToMethod() int {
	if pc == nil {
		return 0
	}
	switch *pc {
	case ALIPAY:
		return 1
	case WECHAT:
		return 2
	case BALANCE:
		return 3
	default:
		return 0
	}
}

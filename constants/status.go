package constants

type OrderStatus int
type OrderStatusString string

const (
	Pending        OrderStatus = 100
	PendingPayment OrderStatus = 200
	PaymentSuccess OrderStatus = 300
	Expire         OrderStatus = 400

	PendingString        OrderStatusString = "pending"
	PendingPaymentString OrderStatusString = "pending_payment"
	PaymentSuccessString OrderStatusString = "payment_success"
	ExpireString         OrderStatusString = "expire"
)

var mapStatusStringToInt = map[OrderStatusString]OrderStatus{
	PendingString:        Pending,
	PendingPaymentString: PendingPayment,
	PaymentSuccessString: PaymentSuccess,
	ExpireString:         Expire,
}

var mapStatusIntToString = map[OrderStatus]OrderStatusString{
	Pending:        PendingString,
	PendingPayment: PendingPaymentString,
	PaymentSuccess: PaymentSuccessString,
	Expire:         ExpireString,
}

func (p OrderStatusString) String() string {
	return string(p)
}
func (p OrderStatus) Int() int {
	return int(p)
}
func (p OrderStatus) GetStatusString() OrderStatusString {
	return mapStatusIntToString[p]
}

func (p OrderStatusString) GetStatusInt() OrderStatus {
	return mapStatusStringToInt[p]
}

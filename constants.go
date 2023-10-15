package exchange_models

type ExchangeName string

type Side string

type OrderType string

type OrderStatus string

var (
	P2PB2B  ExchangeName = "P2PB2B"
	ByBit   ExchangeName = "ByBit"
	Latoken ExchangeName = "Latoken"

	Buy  Side = "Buy"
	Sell Side = "Sell"

	Limit  OrderType = "Limit"
	Market OrderType = "Market"

	Filled            OrderStatus = "Filled"
	PartiallyFilled   OrderStatus = "PartiallyFilled"
	New               OrderStatus = "New"
	Cancelled         OrderStatus = "Cancelled"
	CancelledNotFully OrderStatus = "CancelledNotFully"
)

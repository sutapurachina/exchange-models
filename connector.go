package exchange_models

type Connector interface {
	PostLimitOrder(base, quote string, side Side, baseAmount, price float64, basePrecision, pricePrecision int) (id string, err error)
	CancelOrder(orderId, base, quote string) error
	AllOpenOrders(base, quote string, basePrecision, pricePrecision int) ([]*NetOrder, error)
	OpenOrders(base, quote string, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error)
	OrderBook(base, quote string, side Side, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error)
	FullOrderBook(base, quote string, side Side, basePrecision, pricePrecision int) ([]*NetOrder, error)
	BestBidBestAsk(base, quote string) (bestBid, bestAsk float64, err error)
}

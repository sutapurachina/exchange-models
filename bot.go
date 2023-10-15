package exchange_models

type TradingBot interface {
	PostLimitOrder(order Order) error
	CancelOrder(order Order) error
	GetBestBidAsk() (bestBid float64, bestAsk float64, err error)
}

type BotConfig struct {
	ExName    ExchangeName
	PublicKey string
	SecretKey string
}

type SymbolInfo struct {
	Base           string
	Quote          string
	PricePrecision int
	BasePrecision  int
}

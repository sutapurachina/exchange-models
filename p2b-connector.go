package exchange_models

import (
	"fmt"
	"github.com/sutapurachina/go-p2pb2b"
	"strconv"
)

const (
	P2BBuy  = "buy"
	P2BSell = "sell"
)

type P2BConnector struct {
	Connector
	Client p2pb2b.Client
}

func NewP2BConnector(publicKey, secretKey string) (*P2BConnector, error) {
	client, err := p2pb2b.NewClient(publicKey, secretKey)
	if err != nil {
		return nil, err
	}
	return &P2BConnector{
		Client: client,
	}, nil
}

func (c *P2BConnector) PostLimitOrder(base, quote string, side Side, baseAmount, price float64, basePrecision, pricePrecision int) (id string, err error) {
	orderSide := P2BSell
	if side == Buy {
		orderSide = P2BBuy
	}
	req := &p2pb2b.CreateOrderRequest{
		Market: symbol(base, quote),
		Side:   orderSide,
		Price:  Round(price, pricePrecision),
		Amount: Round(baseAmount, basePrecision),
	}
	resp, err := c.Client.CreateOrder(req)
	if err != nil {
		return
	}
	id = strconv.FormatInt(resp.Result.OrderID, 10)
	return
}

func (c *P2BConnector) CancelOrder(orderId, base, quote string) error {
	numericalOrderId, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return err
	}
	req := &p2pb2b.CancelOrderRequest{
		OrderID: numericalOrderId,
		Market:  symbol(base, quote),
	}
	_, err = c.Client.CancelOrder(req)
	return err
}

func (c *P2BConnector) AllOpenOrders(base, quote string, basePrecision, pricePrecision int) ([]*NetOrder, error) {
	var offset int64 = 0
	res := make([]*NetOrder, 0, 1)
	orders, err := c.OpenOrders(base, quote, basePrecision, pricePrecision, offset, 100)
	if err != nil {
		return nil, err
	}
	res = append(res, orders...)
	for orders != nil && len(orders) > 0 {
		offset += 100
		orders, err = c.OpenOrders(base, quote, basePrecision, pricePrecision, offset, 100)
		res = append(res, orders...)

	}
	return res, nil
}

func (c *P2BConnector) OpenOrders(base, quote string, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error) {
	req := &p2pb2b.QueryUnexecutedRequest{
		Market: symbol(base, quote),
		Offset: offset,
		Limit:  limit,
	}
	resp, err := c.Client.QueryUnexecuted(req)
	if err != nil {
		return nil, err
	}
	res := make([]*NetOrder, 0, 1)
	for _, unexecutedOrder := range resp.Result {
		side := Sell
		if unexecutedOrder.Side == P2BBuy {
			side = Buy
		}
		price, err := strconv.ParseFloat(unexecutedOrder.Price, 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(unexecutedOrder.Amount, 64)
		if err != nil {
			return nil, err
		}
		left, err := strconv.ParseFloat(unexecutedOrder.Left, 64)
		if err != nil {
			return nil, err
		}
		status := New
		if left > 0 {
			status = PartiallyFilled
		}
		orderConfig := &NetOrderConfig{
			Id:           fmt.Sprintf("%d", unexecutedOrder.Id),
			ExName:       P2PB2B,
			Symbol:       symbol(base, quote),
			OrderType:    Limit,
			Side:         side,
			Status:       status,
			Price:        price,
			BaseAmount:   amount,
			BasePrec:     basePrecision,
			PricePrec:    pricePrecision,
			FilledAmount: Round(amount-left, basePrecision),
		}
		order, err := NewNetOrder(orderConfig)
		if err != nil {
			return nil, err
		}
		res = append(res, order)
	}
	return res, err
}

func (c *P2BConnector) OrderBook(base, quote string, side Side, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error) {
	orderSide := P2BBuy
	if side == Sell {
		orderSide = P2BSell
	}
	resp, err := c.Client.GetOrderBook(symbol(base, quote), orderSide, offset, limit)
	if err != nil {
		return nil, err
	}
	orders := make([]*NetOrder, 0, 1)
	for _, order := range resp.Result.Orders {
		side := Sell
		if order.Side == P2BBuy {
			side = Buy
		}
		price := order.Price
		amount := order.Amount
		left := order.Left
		status := New
		if left > 0 {
			status = PartiallyFilled
		}
		orderConfig := &NetOrderConfig{
			Id:           fmt.Sprintf("%d", order.ID),
			ExName:       P2PB2B,
			Symbol:       symbol(base, quote),
			OrderType:    Limit,
			Side:         side,
			Status:       status,
			Price:        price,
			BaseAmount:   amount,
			BasePrec:     basePrecision,
			PricePrec:    pricePrecision,
			FilledAmount: Round(amount-left, basePrecision),
		}
		order, err := NewNetOrder(orderConfig)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (c *P2BConnector) FullOrderBook(base, quote string, side Side, basePrecision, pricePrecision int) ([]*NetOrder, error) {
	var offset int64 = 0
	res := make([]*NetOrder, 0, 1)
	orders, err := c.OrderBook(base, quote, side, basePrecision, pricePrecision, offset, 100)
	if err != nil {
		return nil, err
	}
	res = append(res, orders...)
	for orders != nil && len(orders) > 0 {
		offset += 100
		orders, err = c.OrderBook(base, quote, side, basePrecision, pricePrecision, offset, 100)
		if err != nil {
			return nil, err
		}
		res = append(res, orders...)
	}
	return res, nil
}

func (c *P2BConnector) BestBidBestAsk(base, quote string) (bestBid, bestAsk float64, err error) {
	res, err := c.Client.GetDepthResult(symbol(base, quote), 1)
	if err != nil {
		return 0, 0, err
	}
	bestBid = res.Result.Bids[0][0]
	bestAsk = res.Result.Asks[0][0]
	return
}

func (c *P2BConnector) LastPrice(base, quote string) (lastPrice float64, err error) {
	res, err := c.Client.GetTicker(symbol(base, quote))
	if err != nil {
		return 0, err
	}
	return res.Result.Last, nil
}

func (c *P2BConnector) CurrencyBalance(currency string) (available, freeze float64, err error) {
	req := &p2pb2b.AccountCurrencyBalanceRequest{Currency: currency}
	resp, err := c.Client.PostCurrencyBalance(req)
	if err != nil {
		return
	}
	balance := resp.Result
	available = balance.Available
	freeze = balance.Freeze
	err = nil
	return
}

func (c *P2BConnector) DealHistory(base, quote string, startTime, endTime int64) ([]*Level, error) {
	var offset int64 = 0
	var limit int64 = 100
	req := &p2pb2b.DealsHistoryByMarketRequest{
		Market:    symbol(base, quote),
		StartTime: startTime,
		EndTime:   endTime,
		Offset:    offset,
		Limit:     limit,
	}
	levels := make([]*Level, 0, 1)
	res, err := c.Client.DealsHistoryByMarket(req)
	if err != nil {
		return nil, err
	}
	for _, d := range res.Result.Deals {
		if !d.IsSelfTrade {
			level, err := DealToLevel(d)
			if err != nil {
				return nil, err
			}
			levels = append(levels, level)
		}
	}
	for res.Result.Deals != nil && len(res.Result.Deals) != 0 {
		req.Offset += 100
		res, err = c.Client.DealsHistoryByMarket(req)
		if err != nil {
			return nil, err
		}
		for _, d := range res.Result.Deals {
			if !d.IsSelfTrade {
				level, err := DealToLevel(d)
				if err != nil {
					return nil, err
				}
				levels = append(levels, level)
			}
		}
	}
	return levels, nil
}

func DealToLevel(d p2pb2b.DealHistoryEntry) (*Level, error) {
	price, err := strconv.ParseFloat(d.Price, 64)
	if err != nil {
		return nil, err
	}
	amount, err := strconv.ParseFloat(d.Amount, 64)
	if err != nil {
		return nil, err
	}
	level := &Level{
		Price: price,
	}
	if d.Side == "sell" {
		level.SellAmount = amount
	} else {
		level.BuyAmount = amount
	}
	return level, nil
}

func symbol(base, quote string) string {
	return base + "_" + quote
}

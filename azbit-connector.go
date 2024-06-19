package exchange_models

import (
	azbitgosdk "github.com/sutapurachina/azbit-go-sdk"
	"math"
)

type AzBitConnector struct {
	Connector
	Client *azbitgosdk.AzBitClient
}

func NewAzBitConnector(publicKey, secretKey string) (*AzBitConnector, error) {
	client := azbitgosdk.NewAzBitClient(publicKey, secretKey)

	return &AzBitConnector{
		Client: client,
	}, nil
}

func (c *AzBitConnector) PostLimitOrder(base, quote string, side Side, baseAmount, price float64, basePrecision, pricePrecision int) (id string, err error) {
	orderSide := azbitgosdk.Sell
	if side == Buy {
		orderSide = azbitgosdk.Buy
	}
	return c.Client.PostOrder(orderSide, base, quote, baseAmount, price)
}

func (c *AzBitConnector) CancelOrder(orderId, base, quote string) error {
	return c.Client.CancelOrder(orderId)
}

func (c *AzBitConnector) AllOpenOrders(base, quote string, basePrecision, pricePrecision int) ([]*NetOrder, error) {
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

func (c *AzBitConnector) OpenOrders(base, quote string, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error) {
	orders, err := c.Client.MyOrders(base, quote, "active")
	if err != nil {
		return nil, err
	}
	res := make([]*NetOrder, 0, 1)
	for _, unexecutedOrder := range orders {
		side := Sell
		if unexecutedOrder.IsBid {
			side = Buy
		}
		price := unexecutedOrder.Price
		amount := unexecutedOrder.InitialAmount
		left := unexecutedOrder.Amount
		status := New
		if left > 0 {
			status = PartiallyFilled
		}
		orderConfig := &NetOrderConfig{
			Id:           unexecutedOrder.ID,
			ExName:       AzBit,
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

func (c *AzBitConnector) OrderBook(base, quote string, side Side, basePrecision, pricePrecision int, offset, limit int64) ([]*NetOrder, error) {
	resp, err := c.Client.OrderBook(base, quote)
	if err != nil {
		return nil, err
	}
	orders := make([]*NetOrder, 0, 1)
	for _, order := range resp {
		side := Sell
		if order.IsBid {
			side = Buy
		}
		if side == Sell {
			if order.IsBid {
				continue
			}
		} else {
			if !order.IsBid {
				continue
			}
		}

		price := order.Price
		amount := order.Amount
		left := order.Amount
		status := New
		if left > 0 {
			status = PartiallyFilled
		}
		orderConfig := &NetOrderConfig{
			Id:           "",
			ExName:       AzBit,
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

func (c *AzBitConnector) FullOrderBook(base, quote string, side Side, basePrecision, pricePrecision int) ([]*NetOrder, error) {
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

func (c *AzBitConnector) BestBidBestAsk(base, quote string) (bestBid, bestAsk float64, err error) {
	res, err := c.Client.OrderBook(base, quote)
	if err != nil {
		return 0, 0, err
	}
	bestBid = 0
	bestAsk = math.MaxFloat64
	for _, level := range res {
		if level.IsBid {
			if level.Price > bestBid {
			}
			bestBid = level.Price
		} else {
			if level.Price < bestAsk {
				bestAsk = level.Price
			}
		}
	}
	return
}
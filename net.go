package exchange_models

import (
	"fmt"
	"math"
)

type Net interface {
	Orders(side Side) []Order
	BaseAmount(side Side) float64
	QuoteAmount(side Side) float64
	WeightedAveragePrice(side Side) float64
	Spread() (length, ratio float64)
	InsertOrder(order Order)
	RemoveOrder(id string)

	Print()
}

type ClassicNet struct {
	Net
	BuyOrders  []*NetOrder
	SellOrders []*NetOrder
}

func NewEmptyClassicNet() *ClassicNet {
	buyOrders := make([]*NetOrder, 0, 1)
	sellOrders := make([]*NetOrder, 0, 1)
	return &ClassicNet{
		BuyOrders:  buyOrders,
		SellOrders: sellOrders,
	}
}

func (n *ClassicNet) Orders(side Side) []*NetOrder {
	if side == Buy {
		return n.BuyOrders
	}
	return n.SellOrders
}

func (n *ClassicNet) BaseAmountFromTillLevel(side Side, fromPrice, tillLevel float64) float64 {
	if side == Buy {
		if len(n.BuyOrders) == 0 {
			return 0
		}

		fromPrice, tillLevel = math.Max(fromPrice, tillLevel), math.Min(fromPrice, tillLevel)
		firstIdx := 0
		firstIdxFound := false
		for idx, order := range n.BuyOrders {
			if order.Price() <= fromPrice && !firstIdxFound {
				firstIdx = idx
				firstIdxFound = true
			}
			if order.Price() < tillLevel && firstIdxFound {
				return BaseAmount(n.BuyOrders[firstIdx:idx], n.BuyOrders[0].BasePrecision())
			}
		}
		if firstIdxFound {
			return BaseAmount(n.BuyOrders[firstIdx:], n.BuyOrders[0].BasePrecision())
		}
		return 0
	}

	if len(n.SellOrders) == 0 {
		return 0
	}
	fromPrice, tillLevel = math.Min(fromPrice, tillLevel), math.Max(fromPrice, tillLevel)
	firstIdx := 0
	firstIdxFound := false
	for idx, order := range n.SellOrders {
		if order.Price() >= fromPrice && !firstIdxFound {
			firstIdx = idx
			firstIdxFound = true
		}
		if order.Price() > tillLevel && firstIdxFound {
			return BaseAmount(n.SellOrders[firstIdx:idx], n.SellOrders[0].BasePrecision())
		}
	}
	if firstIdxFound {
		return BaseAmount(n.SellOrders[firstIdx:], n.SellOrders[0].BasePrecision())
	}
	return 0
}

func (n *ClassicNet) BaseAmount(side Side) float64 {
	if side == Buy {
		if len(n.BuyOrders) == 0 {
			return 0
		}
		return BaseAmount(n.BuyOrders, n.BuyOrders[0].BasePrecision())
	}
	if len(n.SellOrders) == 0 {
		return 0
	}
	return BaseAmount(n.SellOrders, n.SellOrders[0].BasePrecision())
}

func (n *ClassicNet) QuoteAmount(side Side) float64 {
	if side == Buy {
		if len(n.BuyOrders) == 0 {
			return 0
		}
		return QuoteAmount(n.BuyOrders)
	}
	if len(n.SellOrders) == 0 {
		return 0
	}
	return QuoteAmount(n.SellOrders)
}

// WeightedAveragePrice todo проходится по списку один раз
func (n *ClassicNet) WeightedAveragePrice(side Side) float64 {
	base := n.BaseAmount(side)
	if base > 0 {
		return n.QuoteAmount(side) / base
	}
	return 0
}

func (n *ClassicNet) Spread() (length, ratio float64) {
	if len(n.BuyOrders) == 0 || len(n.SellOrders) == 0 {
		return 0, 0
	}
	length = n.BuyOrders[0].Price() - n.SellOrders[0].Price()
	ratio = length / n.BuyOrders[0].Price()
	return
}

func (n *ClassicNet) InsertOrder(order *NetOrder) {
	if order.Side() == Buy {
		if len(n.BuyOrders) == 0 {
			n.insertOrderInTheEnd(order)
			return
		}
		for idx, activeOrder := range n.BuyOrders {
			if order.Price() > activeOrder.Price() {
				if idx == 0 {
					n.insertOrderInTheBeginning(order)
					return
				}
				tmp := make([]*NetOrder, idx)
				copy(tmp, n.BuyOrders[0:idx])
				tmp = append(tmp, order)
				n.BuyOrders = append(tmp, n.BuyOrders[idx:]...)
				return
			}
		}
		n.insertOrderInTheEnd(order)
		return
	}
	if len(n.SellOrders) == 0 {
		n.insertOrderInTheEnd(order)
		return
	}
	for idx, activeOrder := range n.SellOrders {
		if order.Price() < activeOrder.Price() {
			if idx == 0 {
				n.insertOrderInTheBeginning(order)
				return
			}
			tmp := make([]*NetOrder, idx)
			copy(tmp, n.SellOrders[0:idx])
			tmp = append(tmp, order)
			n.SellOrders = append(tmp, n.SellOrders[idx:]...)
			return
		}

	}
	n.insertOrderInTheEnd(order)
	return
}

func (n *ClassicNet) insertOrderInTheEnd(order *NetOrder) {
	if order.Side() == Buy {
		n.BuyOrders = append(n.BuyOrders, order)
		return
	}
	n.SellOrders = append(n.SellOrders, order)
}

func (n *ClassicNet) insertOrderInTheBeginning(order *NetOrder) {
	if order.Side() == Buy {
		n.BuyOrders = append([]*NetOrder{order}, n.BuyOrders...)
		return
	}
	n.SellOrders = append([]*NetOrder{order}, n.SellOrders...)
}

func (n *ClassicNet) RemoveOrder(id string) bool {
	side, index, found := n.findOrder(id)
	if !found {
		return false
	}
	n.deleteElementByIdx(side, index)
	return true
}

func (n *ClassicNet) deleteElementByIdx(side Side, idx int) {
	var orders *[]*NetOrder
	if side == Buy {
		orders = &n.BuyOrders
	} else {
		orders = &n.SellOrders
	}
	if idx == 0 {
		*orders = (*orders)[1:]
		return
	}
	ret := make([]*NetOrder, 0, 1)
	ret = append(ret, (*orders)[:idx]...)
	*orders = append(ret, (*orders)[idx+1:]...)
}

func (n *ClassicNet) findOrder(id string) (side Side, index int, found bool) {
	for idx, order := range n.SellOrders {
		if id == order.ID() {
			side = Sell
			index = idx
			found = true
			return
		}
	}
	for idx, order := range n.BuyOrders {
		if id == order.ID() {
			side = Buy
			index = idx
			found = true
			return
		}
	}
	found = false
	return
}

func (n *ClassicNet) Print() {
	fmt.Printf("Net:\n\n\n")
	for i := len(n.SellOrders) - 1; i >= 0; i-- {
		n.SellOrders[i].Print()
	}
	fmt.Println("----------")
	for _, o := range n.BuyOrders {
		o.Print()
	}
}

package exchange_models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type KeyPair struct {
	Public string
	Secret string
	Usdt   float64
	SDFA   float64
}

var (
	dpp1 = KeyPair{
		"509fe4805d01868bb966cee37770cd79",
		"eba845a021e85eb09d25a02b138cfa43",
		7041.4232,
		194.46,
	}

	kz1 = KeyPair{
		"aaac2b64eb0fb8fd12f7c4b854a8b57c",
		"343fdb91a3633fcf2d6293f08fae53bc",
		0,
		0,
	}

	kz2 = KeyPair{
		"803b77b03b7e1e06cea9e3b943052438",
		"d33c5c26133d1eceab59c63fd052c456",
		0,
		900,
	}
	mm = KeyPair{
		"ddfdb39e447b4a3153f6d1150cb7ae9a",
		"c17f06034442ba1ffd473a973d5e5850",
		0,
		0,
	}
)

func TestP2BConnector_GetAllOpenOrders(t *testing.T) {
	c, err := NewP2BConnector(mm.Public, mm.Secret)
	assert.NoError(t, err)
	openOrders, err := c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	for _, order := range openOrders {
		order.Print()
	}
	fmt.Println(len(openOrders))
}

func TestP2BConnector_OrderBook(t *testing.T) {
	c, err := NewP2BConnector(mm.Public, mm.Secret)
	assert.NoError(t, err)
	orders, err := c.FullOrderBook("SDFA", "USDT", Sell, 3, 2)
	assert.NoError(t, err)
	for _, order := range orders {
		order.Print()
	}
	fmt.Println(len(orders))
}

func TestEnemyOrders(t *testing.T) {
	c, err := NewP2BConnector(mm.Public, mm.Secret)
	assert.NoError(t, err)
	buyOrders, err := c.FullOrderBook("SDFA", "USDT", Buy, 3, 2)
	assert.NoError(t, err)
	sellOrders, err := c.FullOrderBook("SDFA", "USDT", Sell, 3, 2)
	assert.NoError(t, err)
	ordersHash := make(map[string]*NetOrder)
	for _, order := range buyOrders {
		ordersHash[order.id] = order
	}
	for _, order := range sellOrders {
		ordersHash[order.id] = order
	}
	ourOrders, err := c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	c, err = NewP2BConnector(kz1.Public, kz1.Secret)
	assert.NoError(t, err)
	orders, err := c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	c, err = NewP2BConnector(kz2.Public, kz2.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	c, err = NewP2BConnector(dpp1.Public, dpp1.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	for _, order := range ourOrders {
		delete(ordersHash, order.ID())
	}

	net := NewEmptyClassicNet()

	for _, order := range ordersHash {
		net.InsertOrder(order)
	}

	net.Print()
	fmt.Println(net.QuoteAmount(Sell))

}

func TestPostCancel(t *testing.T) {
	c, err := NewP2BConnector(mm.Public, mm.Secret)
	assert.NoError(t, err)
	id, err := c.PostLimitOrder("SDFA", "USDT", Buy, 0.001, 53.87, 3, 2)
	assert.NoError(t, err)
	time.Sleep(10 * time.Second)
	assert.NoError(t, c.CancelOrder(id, "SDFA", "USDT"))
}

func TestP2BConnector_BestBidBestAsk(t *testing.T) {
	c, err := NewP2BConnector(mm.Public, mm.Secret)
	assert.NoError(t, err)
	bestBid, bestAsk, err := c.BestBidBestAsk("SDFA", "USDT")
	assert.NoError(t, err)
	fmt.Println(bestBid, bestAsk)
}

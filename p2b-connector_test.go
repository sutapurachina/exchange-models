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
		"153f17bc07ff9cfa1be54f241b1cccdc",
		"a8035fe8f77d9455916ea594d8b3fa79",
		0,
		0,
	}

	o1 = KeyPair{
		"3444084ff5cf5e6ce7b6a33e4c6728b0",
		"ec5ff52b3a86482ed43e3ed605ea3b34",
		9921.8399,
		0,
	}
	mm = KeyPair{
		"ddfdb39e447b4a3153f6d1150cb7ae9a",
		"c17f06034442ba1ffd473a973d5e5850",
		0,
		172.319,
	}
	sp = KeyPair{"139850c9f7727a74fbe1207aeb1f72ca",
		"52326e61af8d68c41898e8876bc2d1fe",
		-300,
		6.185,
	}
	o2 = KeyPair{
		"9bb4ec525c8737e6ac9cc3f635134a48",
		"ea09d3b83fe08f07454243588ad04eb3",
		0,
		0,
	}

	coatTalisman = KeyPair{
		"a00d3950f29b124f3f95c54663cd6339",
		"94e47c3447dcfb0981e40f59fe8471b4",
		0,
		0,
	}

	oeDPP = KeyPair{
		"dc80a68750f6c70b20f64ad4d8a77e90",
		"3fd18671fab23a2a74f02bb32844a429",
		0,
		149.75,
	}
	YDPP = KeyPair{
		"531c3bea70b687d6df5d7c197fb8636d",
		"2929741615b416d24bb3851dbbe9d547",
		0,
		0,
	}
	ADPP = KeyPair{
		"eb20dc5205043121aa501a58973c2f64",
		"095e559fe08ea50c1397b33c19a270ca",
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
	c, err = NewP2BConnector(o1.Public, o1.Secret)
	assert.NoError(t, err)
	orders, err := c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	c, err = NewP2BConnector(o2.Public, o2.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	c, err = NewP2BConnector(dpp1.Public, dpp1.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	c, err = NewP2BConnector(sp.Public, sp.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	for _, order := range ourOrders {
		delete(ordersHash, order.ID())
	}
	c, err = NewP2BConnector(oeDPP.Public, oeDPP.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	for _, order := range ourOrders {
		delete(ordersHash, order.ID())
	}
	c, err = NewP2BConnector(coatTalisman.Public, coatTalisman.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	for _, order := range ourOrders {
		delete(ordersHash, order.ID())
	}
	c, err = NewP2BConnector(YDPP.Public, YDPP.Secret)
	assert.NoError(t, err)
	orders, err = c.AllOpenOrders("SDFA", "USDT", 3, 2)
	assert.NoError(t, err)
	ourOrders = append(ourOrders, orders...)
	for _, order := range ourOrders {
		delete(ordersHash, order.ID())
	}
	c, err = NewP2BConnector(ADPP.Public, ADPP.Secret)
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

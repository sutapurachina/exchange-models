package exchange_models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetOrder_BaseAmount(t *testing.T) {
	side := Sell
	orderConfigs := []*NetOrderConfig{
		{
			Id:         "1",
			Side:       side,
			Price:      4,
			BaseAmount: 1,
		},
		{
			Id:         "2",
			Side:       side,
			Price:      2,
			BaseAmount: 2,
		},
		{
			Id:         "3",
			Side:       side,
			Price:      3,
			BaseAmount: 3,
		},
		{
			Id:         "4",
			Side:       side,
			Price:      1,
			BaseAmount: 4,
		},
		{
			Id:         "5",
			Side:       side,
			Price:      6,
			BaseAmount: 5,
		},
		{
			Id:         "6",
			Side:       side,
			Price:      5.5,
			BaseAmount: 6,
		},
		{
			Id:         "7",
			Side:       side,
			Price:      7,
			BaseAmount: 7,
		},
		{
			Id:         "8",
			Side:       side,
			Price:      9,
			BaseAmount: 8,
		},
		{
			Id:         "9",
			Side:       side,
			Price:      3,
			BaseAmount: 9,
		},
	}
	net := NewEmptyClassicNet()
	for _, oc := range orderConfigs {
		order, err := NewNetOrder(oc)
		assert.NoError(t, err)
		net.InsertOrder(order)
	}

	net.Print()
	net.RemoveOrder("2")
	net.Print()
	net.RemoveOrder("4")
	net.Print()
	net.RemoveOrder("8")
	net.Print()

	fmt.Println(net.BaseAmountFromTillLevel(side, 5, 4))

}

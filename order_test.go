package exchange_models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewNetOrder(t *testing.T) {
	orderConfig := &NetOrderConfig{
		ExName:       P2PB2B,
		Symbol:       "SDFA_USDT",
		Id:           "232323",
		Side:         Buy,
		OrderType:    Limit,
		Price:        33.4,
		BaseAmount:   3,
		FilledAmount: 332,
		PreviousId:   "sdf",
		BasePrec:     3,
		PricePrec:    3,
	}
	order, err := NewNetOrder(orderConfig)
	assert.NoError(t, err)
	order.Print()
}

func TestNetOrder_Marshal(t *testing.T) {
	orderConfig := &NetOrderConfig{
		ExName:       P2PB2B,
		Symbol:       "SDFA_USDT",
		Id:           "232323",
		Side:         Buy,
		OrderType:    Limit,
		Price:        33.4,
		BaseAmount:   3,
		FilledAmount: 332,
		PreviousId:   "sdf",
		BasePrec:     3,
		PricePrec:    3,
	}
	order, err := NewNetOrder(orderConfig)
	assert.NoError(t, err)
	order.Print()
	orderBytes, err := order.Marshal()
	assert.NoError(t, err)
	fmt.Println(string(orderBytes))
	o, err := UnmarshalNetOrder(orderBytes)
	assert.NoError(t, err)
	o.Print()
}

func TestEmptyOrder(t *testing.T) {
	orderConfig := &NetOrderConfig{}
	order, err := NewNetOrder(orderConfig)
	assert.NoError(t, err)
	order.Print()
}

func TestReturn(t *testing.T) {
	orderConfig := &NetOrderConfig{}
	order, err := NewNetOrder(orderConfig)
	assert.NoError(t, err)
	order.Print()

	order.SetCreationDate(time.Now().UTC())
	order.Print()
	order.SetDeathDate(time.Now().UTC().Add(15 * time.Second))
	order.Print()
	order.SetID("133").SetPreviousId("suka").SetStatus(New).Print()
}

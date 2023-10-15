package exchange_models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Order interface {
	ExchangeName() ExchangeName // name of the exchange where order was or is going to be published
	Symbol() string
	ID() string
	SetID(ID string) Order
	Side() Side
	Type() OrderType
	Status() OrderStatus
	SetStatus(status OrderStatus) Order
	Price() float64
	BaseAmount() float64
	QuoteAmount() float64
	FilledQuoteAmount() float64
	FilledAmount() float64
	UnfilledAmount() float64
	AddFilledAmount(amount float64) error
	PreviousId() string //id of order that is logically connected to this order
	SetPreviousId(id string) Order
	CreationDate() time.Time
	SetCreationDate(creationDate time.Time) Order
	DeathDate() time.Time
	SetDeathDate(deathDate time.Time) Order
	Print()
	BasePrecision() int
	PricePrecision() int
	Marshal() ([]byte, error)
}

type NetOrder struct {
	Order
	exchangeName ExchangeName
	symbol       string
	id           string
	side         Side
	orderType    OrderType
	status       OrderStatus
	price        float64
	baseAmount   float64
	filledAmount float64
	previousId   string
	creationDate time.Time
	deathDate    time.Time
	basePrec     int
	pricePrec    int
}

type NetOrderConfig struct {
	ExName       ExchangeName
	Symbol       string
	Id           string
	Side         Side
	OrderType    OrderType
	Status       OrderStatus
	Price        float64
	BaseAmount   float64
	FilledAmount float64
	PreviousId   string
	BasePrec     int
	PricePrec    int
}

func newNetOrder(config *NetOrderConfig) *NetOrder {
	return &NetOrder{
		exchangeName: config.ExName,
		symbol:       config.Symbol,
		id:           config.Id,
		side:         config.Side,
		orderType:    config.OrderType,
		status:       config.Status,
		price:        config.Price,
		baseAmount:   config.BaseAmount,
		filledAmount: config.FilledAmount,
		previousId:   config.PreviousId,
		basePrec:     config.BasePrec,
		pricePrec:    config.PricePrec,
	}
}

// NewNetOrder todo add errors
func NewNetOrder(config *NetOrderConfig) (*NetOrder, error) {
	return newNetOrder(config), nil
}

func (o *NetOrder) ExchangeName() ExchangeName {
	return o.exchangeName
}

func (o *NetOrder) Symbol() string {
	return o.symbol
}

func (o *NetOrder) ID() string {
	return o.id
}

func (o *NetOrder) SetID(ID string) Order {
	o.id = ID
	return o
}

func (o *NetOrder) Side() Side {
	return o.side
}

func (o *NetOrder) Type() OrderType {
	return o.orderType
}

func (o *NetOrder) Status() OrderStatus {
	return o.status
}

func (o *NetOrder) SetStatus(status OrderStatus) Order {
	o.status = status
	return o
}

func (o *NetOrder) Price() float64 {
	return o.price
}

func (o *NetOrder) BaseAmount() float64 {
	return o.baseAmount
}

func (o *NetOrder) QuoteAmount() float64 {
	return o.baseAmount * o.price
}

func (o *NetOrder) FilledQuoteAmount() float64 {
	return o.filledAmount * o.price
}

func (o *NetOrder) FilledAmount() float64 {
	return o.filledAmount
}

func (o *NetOrder) UnfilledAmount() float64 {
	return Round(o.baseAmount-o.filledAmount, o.basePrec)
}

func (o *NetOrder) AddFilledAmount(amount float64) error {
	res := o.filledAmount + amount
	if res > o.baseAmount {
		return errors.New("new base amount is larger than initial")
	}
	o.filledAmount = res
	return nil
}

func (o *NetOrder) PreviousId() string {
	return o.previousId
}

func (o *NetOrder) SetPreviousId(id string) Order {
	o.previousId = id
	return o
}

func (o *NetOrder) CreationDate() time.Time {
	return o.creationDate
}

func (o *NetOrder) SetCreationDate(creationDate time.Time) Order {
	o.creationDate = creationDate
	return o
}

func (o *NetOrder) DeathDate() time.Time {
	return o.deathDate
}

func (o *NetOrder) SetDeathDate(deathDate time.Time) Order {
	o.deathDate = deathDate
	return o
}

func (o *NetOrder) BasePrecision() int {
	return o.basePrec
}

func (o *NetOrder) PricePrecision() int {
	return o.pricePrec
}

func (o *NetOrder) Print() {
	fmt.Printf("%s, %s, %s, %s, %s, price: %f, base amount: %f, filled: %f, previous id: %s, created: %s, ended: %s, base prec: %d, price prec: %d\n",
		o.exchangeName,
		o.symbol,
		o.id,
		o.side,
		o.orderType,
		o.price,
		o.baseAmount,
		o.filledAmount,
		o.previousId,
		o.creationDate,
		o.deathDate,
		o.basePrec,
		o.pricePrec)
}

func (o *NetOrder) Marshal() ([]byte, error) {
	orderBytes, err := json.Marshal(NetOrderConfig{
		ExName:       o.exchangeName,
		Symbol:       o.symbol,
		Id:           o.id,
		Side:         o.side,
		OrderType:    o.orderType,
		Price:        o.price,
		BaseAmount:   o.baseAmount,
		FilledAmount: o.filledAmount,
		PreviousId:   o.previousId,
		BasePrec:     o.basePrec,
		PricePrec:    o.pricePrec,
	})
	if err != nil {
		return nil, err
	}
	return orderBytes, nil
}

func UnmarshalNetOrder(orderBytes []byte) (Order, error) {
	var netOrderConfig NetOrderConfig
	err := json.Unmarshal(orderBytes, &netOrderConfig)
	if err != nil {
		return nil, err
	}
	return NewNetOrder(&netOrderConfig)
}

func BaseAmount(orders []Order, basePrec int) float64 {
	sum := 0.0
	for _, order := range orders {
		sum += order.BaseAmount()
		sum = Round(sum, basePrec)
	}
	return sum
}

func QuoteAmount(orders []Order) float64 {
	sum := 0.0
	for _, order := range orders {
		sum += order.BaseAmount() * order.Price()
	}
	return sum
}

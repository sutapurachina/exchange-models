package exchange_models

import (
	"fmt"
	"testing"
)

var (
	config = &DivideConfig{
		PartsAmount:    5,
		MinPartRatio:   0.15,
		MaxPartRatio:   0.25,
		PricePrecision: 2,
	}
	spread = &Spread{
		TopPrice:    55.12,
		BottomPrice: 54.00,
	}
)

func TestDivide(t *testing.T) {
	res := Divide(config, spread)

	for _, r := range res {
		r.Print()
	}
	fmt.Println("----------------------")
	a := RandomizeSpreadParts(res)
	for _, k := range a {
		k.Print()
	}
}

func TestSplitAmount(t *testing.T) {
	res := SplitAmount(5.834, 3, 3)
	var sum float64
	for _, r := range res {
		fmt.Println(r)
		sum += r
	}
	fmt.Println(sum)
}

package exchange_models

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type DivideConfig struct {
	PartsAmount    int
	MinPartRatio   float64
	MaxPartRatio   float64
	PricePrecision int
}

type Spread struct {
	TopPrice    float64
	BottomPrice float64
}

func (s *Spread) Length() float64 {
	return s.TopPrice - s.BottomPrice
}

func (s *Spread) Print() {
	fmt.Printf("----\n%f\n\n%f\n----\n", s.TopPrice, s.BottomPrice)
}

// Divide делит весь спред на непересекающиеся отрезки
func Divide(config *DivideConfig, baseSpread *Spread) []*Spread {
	parts := make([]*Spread, 0, config.PartsAmount)
	var ratiosSum float64
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bottomPrice := Round(baseSpread.BottomPrice, config.PricePrecision) + math.Pow10(-config.PricePrecision)
	topPrice := Round(baseSpread.TopPrice, config.PricePrecision) - math.Pow10(-config.PricePrecision)
	currentTopPrice := topPrice
	for i := 0; i < config.PartsAmount; i++ {
		ratio := config.MinPartRatio + math.Abs(r.Float64())*(config.MaxPartRatio-config.MinPartRatio)
		ratiosSum += ratio
		nextPrice := Round(topPrice-ratiosSum*baseSpread.Length(), config.PricePrecision)
		if nextPrice < bottomPrice {
			nextPrice = Round(bottomPrice, config.PricePrecision)
		}
		parts = append(parts, &Spread{currentTopPrice, nextPrice})
		currentTopPrice = nextPrice
	}
	return parts
}

func SplitAmount(amount float64, n int, precision int) []float64 {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	amount = Round(amount, precision)
	if n == 1 {
		return []float64{Round(amount, precision)}
	}
	// Generate n random proportions that sum up to 1
	proportions := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		proportions[i] = r.Float64()
		sum += proportions[i]
	}
	for idx, proportion := range proportions {
		proportions[idx] = proportion / sum
	}
	sum = 0
	individualAmounts := make([]float64, n-1, n)
	for i := 0; i < n-1; i++ {
		individualAmounts[i] = Round(proportions[i]*amount, precision)
		sum += individualAmounts[i]
	}
	lastAmount := Round(amount-sum, precision)

	if lastAmount < 0 {
		lastAmount = 0
	}
	individualAmounts = append(individualAmounts, lastAmount)
	return individualAmounts
}

func RandomizeSpreadParts(slice []*Spread) []*Spread {
	// Create a new slice of the same length as the original slice
	randomSlice := make([]*Spread, len(slice))
	// Copy the content of the original slice to the new slice
	copy(randomSlice, slice)

	// Shuffle the new slice randomly using Fisher-Yates algorithm
	for i := len(randomSlice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		randomSlice[i], randomSlice[j] = randomSlice[j], randomSlice[i]
	}

	return randomSlice
}

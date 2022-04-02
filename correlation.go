package main

import (
	"math"
)

func min(a int, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

type Correlation interface {
	Correl() float32
	CorrelTimeX() float32
	CorrelTimeY() float32
}

type Pearson struct {
	x_vals []float64
	y_vals []float64
}

func (p Pearson) Correl() float64 {
	length := min(len(p.x_vals), len(p.y_vals))
	summults := 0.0
	sumXSquares := 0.0
	sumYSquares := 0.0
	SumX := 0.0
	SumY := 0.0
	for i := 0; i < length; i++ {
		x := p.x_vals[i]
		y := p.y_vals[i]
		sumXSquares += math.Pow(x, 2)
		sumYSquares += math.Pow(y, 2)
		summults += x * y
		SumX += x
		SumY += y
	}
	floatlen := float64(length)
	numerator := (floatlen * summults) - (SumX * SumY)
	denominatorPart1 := (floatlen*sumXSquares - math.Pow(SumX, 2))
	denominatorPart2 := (floatlen*sumYSquares - math.Pow(SumY, 2))
	denominator := math.Sqrt(denominatorPart1 * denominatorPart2)
	return numerator / denominator
}

func constrPearson(x_vals []float64, y_vals []float64) Pearson {
	return Pearson{x_vals: x_vals, y_vals: y_vals}
}

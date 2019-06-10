package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
)

const (
	start = 30.0
	min   = 0.0
	max   = 1.0
	vol   = 0.02
)

func genNextPrice(oldPrice float64) float64 {
	rnd := min + rand.Float64()*(max-min)
	changePct := 2 * vol * rnd
	if changePct > vol {
		changePct -= (2 * vol)
	}
	changeAmt := oldPrice * changePct
	newPrice := oldPrice + changeAmt
	return newPrice
}

func getDayPrices() []float64 {
	prices := make([]float64, 96)
	prices[0] = start
	for i := 1; i < 96; i++ {
		prices[i] = genNextPrice(prices[i-1])
	}
	return prices
}

func writeToFile(nums []float64) {
	file, _ := os.Create("prices.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	strings := make([]string, len(nums))

	for i, num := range nums {
		strings[i] = fmt.Sprintf("%v", num)
	}
	writer.Write(strings)
}

func dumpDayPrices() {
	p := getDayPrices()
	writeToFile(p)
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type quote struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

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

func initialData() []byte {
	quotes := []quote{
		{
			Symbol: "HFZ",
			Price:  12.34,
		},
		{
			Symbol: "XYZ",
			Price:  54.55,
		},
	}

	b, err := json.Marshal(quotes)
	if err != nil {
		log.Println(err)
	}

	return b
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

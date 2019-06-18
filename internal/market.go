package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Company struct {
	Id     int
	Symbol string
}

type Quote struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Stock struct {
	symbol string
	price  float64
	min    float64
	max    float64
	vol    float64
}

type Holding struct {
	Uid    int
	Symbol string
	Shares int
}

type Market struct {
	stocks []*Stock
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

func NewMarket() *Market {
	m := new(Market)
	m.stocks = []*Stock{
		{
			symbol: "HFZ",
			price:  11.34,
			min:    0.0,
			max:    1.0,
			vol:    0.02,
		},
		{
			symbol: "IFO",
			price:  6.60,
			min:    0.0,
			max:    1.2,
			vol:    0.03,
		},
	}
	return m
}

func initialData() []byte {
	quotes := []Quote{
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

func (s *Stock) tickPrice() {
	s.price = genNextPrice(s.price)
}

func (m *Market) OpeningBell(broadcast chan []byte) {
	qs := make([]*Quote, len(m.stocks))
	for {
		time.Sleep(time.Second * 2)
		for i, s := range m.stocks {
			s.tickPrice()
			q := &Quote{
				Symbol: s.symbol,
				Price:  s.price,
			}
			qs[i] = q
		}
		qsb, err := json.Marshal(qs)
		if err != nil {
			log.Println(err)
		}
		broadcast <- qsb
	}
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

package gas

import (
	"log"

	"sort"
	"time"
)

type Methods interface {
	SpentPerMonth() (Response, error)
	AveragePerDay() ([]AveragePerDay, error)
	PerHour() (hourly []PricePerHour)
	WholePeriod() (res Value)
}

type Body struct {
	Ethereum `json:"ethereum"`
}

type Ethereum struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Time           string  `json:"time"`
	GasPrice       float64 `json:"gasPrice"`
	GasValue       float64 `json:"gasValue"`
	Average        float64 `json:"average"`
	MaxGasPrice    float64 `json:"maxGasPrice"`
	MedianGasPrice float64 `json:"ice"`
}

func NewTransaction() Methods {
	return &Ethereum{
		Transactions: []Transaction{},
	}
}

// Task #1
// Spent Per Month
type Response struct {
	SpentPerMonth []SpentPerMonth `json:"spentPerMonth"`
}

type SpentPerMonth struct {
	Month    string  `json:"month"`
	GasValue float64 `json:"gasValue"`
}

func (e *Ethereum) SpentPerMonth() (Response, error) {

	perMonths := make([]SpentPerMonth, 0)
	mapMonth := make(map[string]float64, 0)

	temporary := SpentPerMonth{}

	for _, v := range e.Transactions {
		parsed, err := time.Parse("06-01-_2 15:04", v.Time)
		if err != nil {
			log.Println("cannot parse: ", err)
			return Response{nil}, nil
		}

		mapMonth[parsed.Format("2006-01")] += v.GasValue

	}

	keys := make([]string, 0, len(mapMonth))
	for i := range mapMonth {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	for _, v := range keys {
		temporary.Month = v
		temporary.GasValue = mapMonth[v]
		perMonths = append(perMonths, temporary)
	}

	log.Println("SpentperMonth: ", perMonths)

	return Response{perMonths}, nil
}

// Task #2
// Average price per one day
type AveragePerDay struct {
	Day      string  `json:"day"`
	GasPrice float64 `json:"gasPrice"`
}

func (e *Ethereum) AveragePerDay() ([]AveragePerDay, error) {
	daily := make([]AveragePerDay, 0)

	mapDaily := make(map[string]float64, 0)

	temporary := AveragePerDay{}

	for _, v := range e.Transactions {
		parsed, err := time.Parse("06-01-_2 15:04", v.Time)
		if err != nil {
			log.Println("cannot parse: ", err)
			return nil, err
		}

		mapDaily[parsed.Format("2006-01-02")] += v.GasPrice
	}

	keys := make([]string, 0, len(mapDaily))
	for i := range mapDaily {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	for _, v := range keys {
		temporary.Day = v
		temporary.GasPrice = mapDaily[v] / 24

		daily = append(daily, temporary)
	}

	return daily, nil
}

// Task #3
// Частотное распред цены по часам
type PricePerHour struct {
	Hour      string  `json:"Hour"`
	GasPrice  float64 `json:"gasPrice"`
	Frequency float64 `json:"frequency"`
}

func (e *Ethereum) PerHour() (hourly []PricePerHour) {

	period := make(map[string]float64, 0)
	freq := make(map[string]float64, 0)

	// arr := make([]string, 0)

	for _, v := range e.Transactions {
		parsed, err := time.Parse("06-01-_2 15:04", v.Time)
		if err != nil {
			log.Println("cannot parse: ", err)
			return nil
		}
		period[parsed.Format("15:04")] += v.GasPrice
		freq[parsed.Format("15:04")]++

	}

	log.Println("period:", period)

	keys := make([]string, 0, len(period))
	for i := range period {
		keys = append(keys, i)
	}

	sort.Strings(keys)
	res := PricePerHour{}

	// interval := []int{10000, 12500, 15000, 17500, 20000, 25000}
	// ers := Freq(hourly, interval)

	for _, v := range keys {
		res.Hour = v
		res.GasPrice = period[v]
		res.Frequency = freq[v]

		hourly = append(hourly, res)
	}

	return
}

// Task #4
// Value for whole period
type Value struct {
	TotalValue float64 `json:"totalValue"`
}

func (e *Ethereum) WholePeriod() (res Value) {

	arr := make([]float64, 0)

	for _, v := range e.Transactions {
		total := v.GasPrice * v.GasValue
		arr = append(arr, total)
	}

	res.TotalValue = findSum(arr, len(arr))

	return
}

// find sum of each index of arr
func findSum(a []float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	return (findSum(a, n-1) + a[n-1])
}

// function to find freq between interval in arr
func Freq(arr []PricePerHour, interval []int) map[int]float64 {

	if len(interval) <= 0 {
		return nil
	}

	res := make(map[int]float64, 0)
	// max := 0.0

	for _, num := range interval {
		for _, v := range arr {
			// if v.GasPrice > max {
			// 	max = v.GasPrice
			// }
			if v.GasPrice <= float64(num) {
				res[num]++
			}
		}

		// for _, val := range arr {
		// 	switch {
		// 	case val.GasPrice < max && val.GasPrice > max/1.1:
		// 		res[int(math.Round(max))]++
		// 	case val.GasPrice < max/1.1 && val.GasPrice > max/1.2:
		// 		res[int(math.Round(max/1.1))]++
		// 	case val.GasPrice < max/1.2 && val.GasPrice > max/1.3:
		// 		res[int(math.Round(max/1.2))]++
		// 	case val.GasPrice < max/1.3 && val.GasPrice > max/1.4:
		// 		res[int(math.Round(max/1.3))]++
		// 	case val.GasPrice < max/1.4 && val.GasPrice > max/1.5:
		// 		res[int(math.Round(max/1.4))]++
		// 	case val.GasPrice < max/1.5 && val.GasPrice > max/1.7:
		// 		res[int(math.Round(max/1.5))]++
		// 	case val.GasPrice < max/1.7 && val.GasPrice >= max/2:
		// 		res[int(math.Round(max/1.7))]++
		// 	case val.GasPrice < max/2 && val.GasPrice >= max/2.1:
		// 		res[int(math.Round(max/1.8))]++
		// 	case val.GasPrice < max/2. && val.GasPrice >= max/2.5:
		// 		res[int(math.Round(max/2))]++
		// 	}
	}

	return res
}

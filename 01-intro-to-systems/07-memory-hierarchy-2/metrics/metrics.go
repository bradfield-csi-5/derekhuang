package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserId int

type DollarAmount struct {
	dollars uint64
    cents float64
}

func AverageAge(ages []float64) float64 {
	average := 0.0
    count := len(ages)
    var i int
    for i = 0; i < count - 3; i += 4 {
        average += ages[i] + ages[i+1] + ages[i+2] + ages[i+3]
    }
	for ; i < count; i++ {
		average += ages[i]
	}
	return average / float64(count)
}

func AveragePaymentAmount(payments []DollarAmount) float64 {
	average := 0.0
    count := len(payments)
    var i int
    for i = 0; i < count - 3; i += 4 {
        average += float64(payments[i].dollars) + payments[i].cents
        average += float64(payments[i+1].dollars) + payments[i+1].cents
        average += float64(payments[i+2].dollars) + payments[i+2].cents
        average += float64(payments[i+3].dollars) + payments[i+3].cents
    }
	for ; i < count; i++ {
		average += float64(payments[i].dollars) + payments[i].cents
	}
	return average / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []DollarAmount) float64 {
	avg := AveragePaymentAmount(payments)
	squaredDiffs := 0.0
    count := float64(len(payments))
    for _, p := range payments {
        amount := float64(p.dollars) + p.cents
        diff := amount - avg
        squaredDiffs += diff * diff
    }
	return math.Sqrt(squaredDiffs / count)
}

func StdDevPaymentAmountAltVariance(payments []DollarAmount) float64 {
    avg := AveragePaymentAmount(payments)
    total := 0.0
    count := float64(len(payments))
    for _, p := range payments {
        amount := float64(p.dollars) + p.cents
        total += amount * amount
    }
    return math.Sqrt((total / count) - (avg * avg))
}

func LoadData() ([]float64, []DollarAmount) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

    var ages []float64
	for _, line := range userLines {
		age, _ := strconv.ParseFloat(line[2], 64)
        ages = append(ages, age)
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

    var payments []DollarAmount
	for _, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
        payments = append(payments, DollarAmount{
            uint64(paymentCents / 100),
            float64(paymentCents % 100) / 100,
        })
	}

    return ages, payments
}

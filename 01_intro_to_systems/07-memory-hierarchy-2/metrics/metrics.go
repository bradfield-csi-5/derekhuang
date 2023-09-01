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
	dollars, cents uint64
}

type Payment struct {
	amount DollarAmount
}

func AverageAge(userAges []float64) float64 {
	average, count := 0.0, 0.0
	for _, age := range userAges {
		count += 1
		average += age
	}
	return average / count
}

func AveragePaymentAmount(payments []Payment) float64 {
	average, count, total_cents := 0.0, 0.0, 0.0
    for _, p := range payments {
        count += 1
        amount := float64(p.amount.dollars)
        total_cents += float64(p.amount.cents)
        average += amount
    }
	return average / count + (total_cents / (count * 100))
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []Payment) float64 {
	mean := AveragePaymentAmount(payments)
	squaredDiffs, count := 0.0, 0.0
    for _, p := range payments {
        count += 1
        amount := float64(p.amount.dollars) + float64(p.amount.cents)/100
        diff := amount - mean
        squaredDiffs += diff * diff
    }
	return math.Sqrt(squaredDiffs / count)
}

func LoadData() ([]float64, []Payment) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

    var userAges []float64
	for _, line := range userLines {
		// age, _ := strconv.Atoi(line[2])
		age, _ := strconv.ParseFloat(line[2], 64)
        userAges = append(userAges, age)
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

    var payments []Payment
	for _, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
        payments = append(payments, Payment{
			DollarAmount{uint64(paymentCents / 100), uint64(paymentCents % 100)},
		})
	}

    return userAges, payments
}

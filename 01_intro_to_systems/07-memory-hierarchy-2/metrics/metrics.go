package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserId int
type UserMap map[UserId]*User

type DollarAmount struct {
	dollars, cents uint64
}

type Payment struct {
	amount DollarAmount
}

type User struct {
	id       UserId
	age      int
}

func AverageAge(users []User) float64 {
	average, count := 0.0, 0.0
	for _, u := range users {
		count += 1
		average += float64(u.age)
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

func LoadData() ([]User, []Payment) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

    var users []User
	for _, line := range userLines {
		id, _ := strconv.Atoi(line[0])
		age, _ := strconv.Atoi(line[2])
        users = append(users, User{UserId(id), age})
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

    return users, payments
}

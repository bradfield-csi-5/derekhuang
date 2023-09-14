package metrics

import (
	"math"
	"testing"
)

func BenchmarkMetrics(b *testing.B) {
	ages, payments := LoadData()

	b.Run("Average age", func(b *testing.B) {
		actual := 0.0
		for n := 0; n < b.N; n++ {
			actual = AverageAge(ages)
		}
		expected := 59.62
		if math.IsNaN(actual) || math.Abs(actual-expected) > 0.01 {
			b.Fatalf("Expected average age to be around %.2f, not %.3f", expected, actual)
		}
	})

	b.Run("Average payment", func(b *testing.B) {
		actual := 0.0
		for n := 0; n < b.N; n++ {
			actual = AveragePaymentAmount(payments)
		}

		expected := 499850.559
		if math.IsNaN(actual) || math.Abs(actual-expected) > 0.01 {
			b.Fatalf("Expected average payment amount to be around %.2f, not %.3f", expected, actual)
		}
	})

	b.Run("Payment stddev", func(b *testing.B) {
		actual := 0.0
		for n := 0; n < b.N; n++ {
			actual = StdDevPaymentAmount(payments)
		}
		expected := 288684.850
		if math.IsNaN(actual) || math.Abs(actual-expected) > 0.01 {
			b.Fatalf("Expected standard deviation to be around %.2f, not %.3f", expected, actual)
		}
	})

	b.Run("Payment stddev alt", func(b *testing.B) {
		actual := 0.0
		for n := 0; n < b.N; n++ {
			actual = StdDevPaymentAmountAltVariance(payments)
		}
		expected := 288684.850
		if math.IsNaN(actual) || math.Abs(actual-expected) > 0.01 {
			b.Fatalf("Expected standard deviation to be around %.2f, not %.3f", expected, actual)
		}
	})

}
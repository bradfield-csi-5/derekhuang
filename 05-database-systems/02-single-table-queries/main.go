package main

import "fmt"

var movies = []Row{
	{
		Values: []Value{
			{"id", "4995"},
			{"title", "Beautiful Mind, A (2001)"},
			{"genre", "Drama|Romance"},
		},
	},
	{
		Values: []Value{
			{"id", "4996"},
			{"title", "Little Otik (Otesanek) (2000)"},
			{"genre", "Comedy|Drama|Fantasy"},
		},
	},
	{
		Values: []Value{
			{"id", "4997"},
			{"title", "Convent, The (2000)"},
			{"genre", "Horror|Sci-Fi"},
		},
	},
	{
		Values: []Value{
			{"id", "4998"},
			{"title", "Defiant Ones, The (1958)"},
			{"genre", "Adventure|Crime|Drama|Thriller"},
		},
	},
	{
		Values: []Value{
			{"id", "4999"},
			{"title", "Dodsworth (1936)"},
			{"genre", "Drama|Romance"},
		},
	},
	{
		Values: []Value{
			{"id", "5000"},
			{"title", "Medium Cool (1969)"},
			{"genre", "Drama|Romance"},
		},
	},
	{
		Values: []Value{
			{"id", "5001"},
			{"title", "Sahara (1943)"},
			{"genre", "Action|Drama|War"},
		},
	},
	{
		Values: []Value{
			{"id", "5002"},
			{"title", "Fritz the Cat (1972)"},
			{"genre", "Animation"},
		},
	},
	{
		Values: []Value{
			{"id", "5003"},
			{"title", "Nine Lives of Fritz the Cat, The (1974)"},
			{"genre", "Animation"},
		},
	},
	{
		Values: []Value{
			{"id", "5004"},
			{"title", "Party, The (1968)"},
			{"genre", "Comedy"},
		},
	},
	{
		Values: []Value{
			{"id", "5005"},
			{"title", "Separate Tables (1958)"},
			{"genre", "Drama"},
		},
	},
}

type Value struct {
	Name   string
	StrVal string
}

type Row struct {
	Values []Value
}

type Operator interface {
	Next() bool
	Execute() Row
}

type ScanOperator struct {
	rows []Row
	idx  int
}

func (s *ScanOperator) Next() bool {
	s.idx++
	return s.idx < len(s.rows)
}

func (s *ScanOperator) Execute() Row {
	return s.rows[s.idx]
}

func NewScanOperator(rows []Row) Operator {
	return &ScanOperator{
		rows: rows,
		idx:  -1,
	}
}

func main() {
	scanner := NewScanOperator(movies)
	for scanner.Next() {
		fmt.Println(scanner.Execute())
	}
}

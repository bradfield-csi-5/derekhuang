/*
Exercise 4.10: Modify issues to report the results in age categories, say less than a month old, less than a year old, and more than a year old.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"gopl.io/ch4/github"
)

func main() {
	now := time.Now()
	var monthAgo time.Time = now.AddDate(0, -1, 0)
	var yearAgo time.Time = now.AddDate(-1, 0, 0)
	var ltMonth []*github.Issue
	var ltYear []*github.Issue
	var gtYear []*github.Issue
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].CreatedAt.After(result.Items[j].CreatedAt)
	})
	for _, item := range result.Items {
		if item.CreatedAt.After(monthAgo) {
			ltMonth = append(ltMonth, item)
		} else if item.CreatedAt.After(yearAgo) {
			ltYear = append(ltYear, item)
		} else {
			gtYear = append(gtYear, item)
		}
	}
	fmt.Printf("%d issues\n\n", result.TotalCount)
	fmt.Print("Issues created within the last month:\n")
	for _, item := range ltMonth {
		fmt.Printf("#%-7d\t%v\t%9.9s\t%.55s\n",
			item.Number, item.CreatedAt.Format(time.DateOnly), item.User.Login, item.Title)
	}
	fmt.Print("\n")
	fmt.Print("Issues created within the last year:\n")
	for _, item := range ltYear {
		fmt.Printf("#%-7d\t%v\t%9.9s\t%.55s\n",
			item.Number, item.CreatedAt.Format(time.DateOnly), item.User.Login, item.Title)
	}
	fmt.Print("\n")
	fmt.Print("Issues created more than a year ago:\n")
	for _, item := range gtYear {
		fmt.Printf("#%-7d\t%v\t%9.9s\t%.55s\n",
			item.Number, item.CreatedAt.Format(time.DateOnly), item.User.Login, item.Title)
	}
}

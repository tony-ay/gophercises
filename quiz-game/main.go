package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	file      *string
	timeLimit *int
)

func init() {
	file = flag.String("file", "./problems.csv", "path to quiz csv file")
	timeLimit = flag.Int("limit", 30, "time limit for quiz")
}

func main() {
	flag.Parse()

	f, err := os.Open(*file)

	if err != nil {
		fmt.Printf("Failed to open file: %s\n", *file)
		os.Exit(1)
	}

	lines, err := csv.NewReader(f).ReadAll()

	if err != nil {
		fmt.Printf("Failed to parse %s as CSV file\n", *file)
		os.Exit(1)
	}

	fmt.Println("Welcome to the quiz. Answer each question. Press enter to begin.")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	done := make(chan int)

	var numCorrect int

	go func() {
		defer close(done)
		for i, l := range lines {
			q := l[0]
			a, err := strconv.Atoi(l[1])

			if err != nil {
				fmt.Printf("Failed to parse answer \"%s\" for line %d\n", l[1], i)
				os.Exit(1)
			}

			var ans int
			fmt.Printf("What is %s?\n", q)
			fmt.Scanf("%d", &ans)

			if ans == a {
				fmt.Println("Correct.")
				numCorrect += 1
			} else {
				fmt.Println("Incorrect.")
			}
		}
	}()

	select {
	case <-timer.C:
		fmt.Println("Timer ran out.")
	case <-done:
		fmt.Println("Quiz finished.")
	}
	fmt.Printf("%d out of %d correct.\n", numCorrect, len(lines))
}

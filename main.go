package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)
	shuffleProblems(&problems)

	score := runQuize(problems, timeLimit)

	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}

func runQuize(problems []problem, timeLimit *int) int {
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	score := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			return score
		case answer := <-answerCh:
			if answer == p.answer {
				score++
			}
		}
	}
	return score
}

func shuffleProblems(problems *[]problem) {
	rand.Shuffle(len(*problems), func(i, j int) {
		(*problems)[i], (*problems)[j] = (*problems)[j], (*problems)[i]
	})
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

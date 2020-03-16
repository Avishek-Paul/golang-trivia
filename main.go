package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type problem struct {
	question string
	options  [4]string
	answer   int
}

// define flags as global variables
var (
	csvFilePath string
	randomize   bool
	quizTime    int
)

func init() {
	flag.StringVar(&csvFilePath, "csv", "questions.csv", "path to csv file containg questions")
	flag.BoolVar(&randomize, "random", false, "if included, randomizes question order")
	flag.IntVar(&quizTime, "time", 10, "duration of quiz in seconds")
	flag.Parse()
}

func main() {

	//  Read in the problems
	problems := readProblemsCSV()
	numCorrect := 0

	fmt.Println("---Generic Quiz Game---")
	fmt.Println("Press [Enter] to start")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	timeout := time.After(time.Second * time.Duration(quizTime))
	ansChannel := make(chan int)

	// loop through the problems
loop:
	for i, problem := range problems {

		// ask the questions and display the options
		fmt.Printf("%d. %s\n", i+1, problem.question)
		for j, choice := range problem.options {
			fmt.Printf("\t%d. %s\n", j+1, choice)
		}
		// ask for user input asynchronously using goroutine
		//passing in ansChannel which will contain answer
		go func(outAnsChannel chan int) {
			fmt.Printf("Enter your answer: ")
			var userInput int
			fmt.Scanf("%d", &userInput)
			outAnsChannel <- userInput
		}(ansChannel)

		// select statement to check time and if answer received
		select {
		case <-timeout:
			fmt.Println("Time is up!")
			break loop
		case ans, _ := <-ansChannel:
			if ans == problem.answer {
				numCorrect++
			}
		}
	}

	fmt.Printf("You received a final score of %.2f%%\n", float64(numCorrect)/float64(len(problems))*100)
}

func readProblemsCSV() []problem {

	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatalf("Could not open file at path: %s\nError:\n\t'%s'", csvFilePath, err)
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	var problems []problem = make([]problem, len(lines))

	for idx, line := range lines {
		var currOptions [4]string
		copy(currOptions[:], line[1:5])
		currAnswer, _ := strconv.ParseInt(line[5], 10, 3)
		problems[idx] = problem{
			question: line[0],
			options:  currOptions,
			answer:   int(currAnswer),
		}
	}

	if randomize {
		problems = shuffleProblems(problems)
	}

	return problems
}

func shuffleProblems(problems []problem) []problem {
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
	return problems
}

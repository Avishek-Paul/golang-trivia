package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type problem struct {
	question string
	options  [4]string
	answer   int
}

func main() {

	//  Read in the problems
	problems := readProblemsCSV("questions.csv")
	numCorrect := 0

	fmt.Println("---Generic Quiz Game---")
	fmt.Println("Press [Enter] to start")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	timeout := time.After(time.Second * time.Duration(5))
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

func readProblemsCSV(filepath string) []problem {
	csvFilename := flag.String("csv", filepath, "a csv file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
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
	return problems
}

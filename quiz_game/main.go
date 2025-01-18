package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const csv_path = "questions_answers.csv"
const quiz_timer_duration = 20

func ReadCSVFile(filepath string) [][]string {

	file, err := os.Open(filepath)
	if err != nil {
		file.Close()
		return make([][]string, 0)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		file.Close()
		return make([][]string, 0)
	}

	file.Close()
	return records
}

func RunQuiz(questions_answers [][]string) int {

	correct_answers := 0
	timer := time.NewTicker(quiz_timer_duration * time.Second)
	quiz_done := make(chan bool)

	go func() {
		for _, question_answer := range questions_answers {
			fmt.Println(question_answer[0])

			answer_chan := make(chan string)

			go func() {
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				answer_chan <- scanner.Text()
			}()

			select {
			case <-timer.C:
				fmt.Println("TIMES UP")
				quiz_done <- true
				return
			case text := <-answer_chan:
				answer, err := strconv.ParseInt(text, 10, 64)
				if len(text) == 0 || err != nil {
					fmt.Println("WRONG!")
					continue
				}

				expected_answer, err := strconv.ParseInt(strings.TrimSpace(question_answer[1]), 10, 64)
				if err != nil {
					fmt.Println("Something went wrong on our end")
					continue
				}

				if answer == expected_answer {
					fmt.Println("CORRECT")
					correct_answers++
				} else {
					fmt.Println("WRONG! ERRR")
				}
			}
		}
		quiz_done <- true
	}()

	<-quiz_done
	timer.Stop()

	return correct_answers
}

func main() {
	records := ReadCSVFile(csv_path)

	if len(records) == 0 {
		fmt.Println("Can't load the quiz :(")
		return
	}

	correct_answers := RunQuiz(records)

	var builder strings.Builder
	builder.WriteString("You got ")
	builder.WriteString(fmt.Sprintf("%d", correct_answers))

	if correct_answers == 1 {
		builder.WriteString(" question correct!")
	} else {
		builder.WriteString(" questions correct!")
	}

	fmt.Println(builder.String())
}

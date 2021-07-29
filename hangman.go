package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Hangman struct {
	Attempts   uint8
	WordToFind string
	ToRev      []int
	StockChar  []string
}

// isInside returns true if a value (int) is at least one time in a aray (of int) othertwise it returns false
func isInside(value *int, arr *[]int) bool {
	for _, elem := range *arr {
		if elem == *value {
			return true
		}
	}
	return false
}

// isInside returns true if a value (int) is at least one time in a aray (of int) othertwise it returns false
func isInsideChar(value *string, s *string) bool {
	for _, elem := range *s {
		if string(elem) == *value {
			return true
		}
	}
	return false
}

// printHangman prints to console the hangman at a given status
func printHangman(hangman []string, status *uint8) {
	for i := ((*status - 1) * 7); i < ((*status-1)*7)+7; i++ {
		fmt.Println(hangman[i])
	}
	fmt.Println()
}

// printWordProgress prints the progess of finding the word
func printWordProgress(wordToFind *string, toRev *[]int) {
	wordToFindLen := len(*wordToFind) - 1

	for index, char := range *wordToFind {
		if isInside(&index, toRev) {
			fmt.Print(strings.ToUpper(string(char)))
		} else {
			fmt.Print("_")
		}

		if index != wordToFindLen {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println()
}

// checkError checks if the error is different from nil otherwise displays error
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

//readFile returns an array of string which is the same as the file (line = line) without useless lines
func readFile(Filename string) []string {
	var source []string
	file, err := os.Open(Filename) // opens the .txt
	checkError(err)
	scanner := bufio.NewScanner(file) // scanner scans the file
	scanner.Split(bufio.ScanLines)    // sets-up scanner preference to read the file line-by-line
	for scanner.Scan() {              // loop that performs a line-by-line scan on each new iteration
		if scanner.Text() != "" {
			source = append(source, scanner.Text()) // adds the value of scanner (that contains the characters from StylizedFile) to source
		}
	}
	file.Close() // closes the file
	return source
}

// Open file to get random word
func randomWord() string {
	words := readFile(os.Args[1])
	return strings.ToUpper(words[rand.Intn(len(words))])
}

func testWord(status *Hangman, UserTry *string) (valid bool, elem string) {
	// If the user entry is a char
	if len(*UserTry) <= 1 {

		AllChar := strings.Join(status.StockChar, "")
		// Add the char to StockChar
		if !isInsideChar(UserTry, &AllChar) {
			status.StockChar = append(status.StockChar, *UserTry)
		} else {
			fmt.Println("Already try", *UserTry)
			return false, "char"
		}

		if isInsideChar(UserTry, &status.WordToFind) {
			for index, char := range status.WordToFind {
				if string(char) == *UserTry {
					if !isInside(&index, &status.ToRev) {
						status.ToRev = append(status.ToRev, index)
					}
				}
			}
			return true, "char"
		} else {
			status.Attempts++
			return false, "char"
		}
	} else { //If it's a word
		if *UserTry == status.WordToFind {
			return true, "word"
		} else {
			status.Attempts += 2
			return false, "word"
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Check if user provide file containing words
	if len(os.Args[1:]) <= 0 {
		fmt.Println("Missing files of words.")
		return
	}

	// Get hangman
	hangman := readFile("hangman.txt")

	// Retrieve the chosen random word
	wordToFind := randomWord()

	// Number of letter to reveal
	reveal := len(wordToFind)/2 - 1

	var toRev []int
	if reveal > 0 {
		var randInt int

		for i := 0; i < reveal; i++ {
			randInt = rand.Intn(len(wordToFind))

			if !isInside(&randInt, &toRev) {
				toRev = append(toRev, randInt)
			} else {
				i--
			}
		}
	}

	status := Hangman{Attempts: 0, WordToFind: wordToFind, ToRev: toRev, StockChar: []string{}}

	fmt.Println("Good luck, you have 10 attempts.")
	printWordProgress(&wordToFind, &toRev)

	var UserTry string

	for status.Attempts != 10 {
		// Get user input
		fmt.Print("Choose: ")
		fmt.Scanln(&UserTry)
		UserTry = strings.ToUpper(UserTry)

		// Test if user input match something in the word
		valid, elem := testWord(&status, &UserTry)

		if valid && elem == "word" {
			fmt.Println("Congrats !")
			break
		} else if valid && elem == "char" {
			printWordProgress(&status.WordToFind, &status.ToRev)

			if len(status.ToRev) == len(status.WordToFind) {
				fmt.Println("Congrats !")
				break
			}
		} else {
			fmt.Println("Not present in the word,", 10-status.Attempts, "attempts remaining")
			if status.Attempts > 10 {
				status.Attempts = 10
			}
			printHangman(hangman, &status.Attempts)
		}
	}

	// If the maximum number of try reached
	if status.Attempts == 10 {
		fmt.Println("Failed")
		fmt.Println("Word to find was", status.WordToFind)
	}
}

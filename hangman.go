package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// isInside returns true if a value (int) is at least one time in a aray (of int) othertwise it returns false
func isInside(value *int, arr *[]int) bool {
	for _, elem := range *arr {
		if elem == *value {
			return true
		}
	}
	return false
}

// printHangman prints to console the hangman at a given status
func printHangman(hangman []string, status *uint8) {
	for i := (*status * 7); i < (*status*7)+7; i++ {
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

func main() {
	// Check if user provide file containing words
	if len(os.Args[1:]) <= 0 {
		fmt.Println("Missing files of words.")
		return
	}

	// Get hangman
	hangman := readFile("hangman.txt")

	// Open file to get word
	rand.Seed(time.Now().UnixNano())
	words := readFile(os.Args[1])
	wordToFind := strings.ToUpper(words[rand.Intn(len(words))])

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

	fmt.Println("Good luck, you have 10 attempts.")
	printWordProgress(&wordToFind, &toRev)

	var userInput string
	var attempts uint8
	var charInside bool

	for attempts != 10 {
		charInside = false
		// Get user input
		fmt.Print("Choose: ")
		fmt.Scanln(&userInput)

		userInput = strings.ToUpper(userInput)
		if len(userInput) > 1 {
			if userInput == wordToFind {
				fmt.Println("Congrats !")
				break
			} else {
				attempts += 2
			}
		}

		// Check is choosen letter is in the word
		for index, char := range wordToFind {
			if strings.EqualFold(string(char), userInput) {
				if !isInside(&index, &toRev) {
					toRev = append(toRev, index)
				}
				charInside = true
			}
		}

		// Si lettre proposé non pas dans mot
		if !charInside {
			fmt.Println("Not present in the word,", 10-attempts-1, "attempts remaining")

			// Display Hangman
			printHangman(hangman, &attempts)
			attempts++
		} else {
			// Print word progess
			printWordProgress(&wordToFind, &toRev)
		}

		// If word found
		if len(toRev) == len(wordToFind) {
			fmt.Println("Congrats !")
			break
		}
	}

	// Si nombre d'essai max atteint
	if attempts == 10 {
		fmt.Println("Failed")
	}
}

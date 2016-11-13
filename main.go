package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	maxLoopCount   = 10
	wordsPerMinute = 500
)

const (
	clear = "\033[2J"

	// Position Cursor (argument order: Line, Column)
	positionCursor = "\033[%d;%dH"
)

func getInput() string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		log.Fatal("no input")
	}

	b := make([]byte, 1024*4)
	n, err := os.Stdin.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return string(b[:n])
}

// Using tput, determines the dimensions of the users terminal so we
// have a canvas to draw on
func termDimensions() (int, int, error) {

	// Grab the terminal dimensions
	linesCmd := exec.Command("tput", "lines")
	colsCmd := exec.Command("tput", "cols")

	lines, err := linesCmd.Output()
	if err != nil {
		return 0, 0, err
	}

	cols, err := colsCmd.Output()
	if err != nil {
		return 0, 0, err
	}

	lines = bytes.TrimSpace(lines)
	cols = bytes.TrimSpace(cols)

	linesInt, err := strconv.Atoi(string(lines))
	if err != nil {
		return 0, 0, err
	}

	colsInt, err := strconv.Atoi(string(cols))
	if err != nil {
		return 0, 0, err
	}

	return linesInt, colsInt, nil
}

func main() {

	// Take input from the terminal
	input := getInput()
	words := strings.Fields(input)

	wordsLen := len(words)
	wordsIndex := 0

	// @TODO(mark): Convert to flag, termDimensions doesn't use caller tty
	rows, cols, err := termDimensions()
	if err != nil {
		log.Fatal(err)
	}
	rows = 37
	cols = 176

	middleCol := cols / 2
	middleRow := rows / 2

	for {
		// Allows this to repeat forever
		loop := wordsIndex / wordsLen
		wordIndex := wordsIndex % wordsLen

		if loop > maxLoopCount {
			break
		}

		word := words[wordIndex]
		col := middleCol - (len(word) / 2)
		fmt.Print(clear +
			fmt.Sprintf(positionCursor, middleRow, col) +
			word +
			fmt.Sprintf(positionCursor, rows, cols))

		wordsIndex++

		// @TODO(mark): Alter the sleep based on the length of the word
		time.Sleep(time.Minute / wordsPerMinute)
	}
}

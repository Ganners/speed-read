package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	// Default lines and columns. We need to know these to work out how to
	// center aline the text to the terminal
	defaultLines = 30
	defaultCols  = 80

	// Default number of times to loop and speed to play at
	defaultLoopCount      = 1
	defaultWordsPerMinute = 200

	// Help messages
	helpLines = "lines in terminal, e.g. '$(tput lines)'"
	helpCols  = "lines in terminal, e.g. '$(tput cols)'"

	helpWordsPerMinute = "words per minute, recommend somewhere from 300 - 700"
	helpLoopCount      = "number of times to loop, 0 for infinite"

	// The word length which we consider we can read in a normal time
	// The average word length is actually 5.1 letters (interesting fact!)
	averageWordLength = 5

	// For each letter over the average, what should the percentage increase in
	// pause be
	percentageIncreasePerLetter = 20
	percentageIncreasePause     = 50

	// Escape sequences to clear the terminal and position the cursor
	clear          = "\033[2J"
	positionCursor = "\033[%d;%dH"
)

func main() {
	// Take input from the terminal
	input, err := getInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Parses flags etc.
	opts := getOptions()
	middleCol := opts.cols / 2
	middleRow := opts.lines / 2

	// Convert to individual words
	words := strings.Fields(input)
	wordsLen := uint(len(words))
	wordsIndex := uint(0)

	for {
		// Allows this to repeat forever
		loop := wordsIndex / wordsLen
		wordIndex := wordsIndex % wordsLen

		if opts.loop > 0 && loop >= opts.loop {
			break
		}

		word := words[wordIndex]
		wordLen := uint(len(word))
		col := middleCol - (wordLen / 2)

		// Construct and print the frame
		frame := bytes.NewBuffer(make([]byte, 256))
		{
			// Clear the terminal
			frame.WriteString(clear)

			// Position the cursor
			frame.WriteString(fmt.Sprintf(positionCursor, middleRow, col))

			// Print the word
			frame.WriteString(word)

			// Position the cursor 'out of sight'
			frame.WriteString(fmt.Sprintf(positionCursor, opts.lines, opts.cols))

			// Print the frame
			fmt.Print(frame.String())
		}

		wordsIndex++

		// The sleep time will in general be the speed specified by the user.
		// However, it will be scaled up for words that are longer than the
		// average so there is a bit more time to read them
		sleepTime := time.Minute / time.Duration(opts.wpm)
		if wordLen > averageWordLength {
			increase := 100 + ((wordLen - averageWordLength) * percentageIncreasePerLetter)
			sleepTime = (sleepTime / 100) * time.Duration(increase)
		}

		// If there should be a pause
		if shouldPause(word) {
			sleepTime = (sleepTime / 100) * (100 + time.Duration(percentageIncreasePause))
		}

		time.Sleep(sleepTime)
	}
}

// The user options
type options struct {
	lines uint
	cols  uint
	wpm   uint
	loop  uint
}

// Retrieves user input from flags
func getOptions() *options {
	opts := options{}
	flag.UintVar(&opts.lines, "lines", defaultLines, helpLines)
	flag.UintVar(&opts.cols, "cols", defaultCols, helpCols)
	flag.UintVar(&opts.wpm, "wpm", defaultWordsPerMinute, helpWordsPerMinute)
	flag.UintVar(&opts.loop, "loop", defaultLoopCount, helpLoopCount)

	flag.Parse()

	return &opts
}

// Retrieves the user input that has been piped in from stdin, e.g.
// echo "this is some input" | ./speed-reading
func getInput() (string, error) {
	// First check if there is some input
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("could not check stat: %s", err.Error())
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return "", errors.New("no input")
	}

	// Allocate some bytes and read in input
	b := make([]byte, 1024*1024*1)
	n, err := os.Stdin.Read(b)
	if err != nil {
		return "", fmt.Errorf("could not read from stdin: %s", err.Error())
	}

	return string(b[:n]), nil
}

// For a given word, works out if there should be an additional pause
func shouldPause(word string) bool {
	switch word[len(word)-1] {
	case '.', ',', ':':
		return true
	}
	return false
}

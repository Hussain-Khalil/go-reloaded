package main

import (
	"fmt"
	"os"
	"strings"

	"goreloaded/processor"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt output.txt")
		return
	}

	chars, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error")
		return
	}

	text := string(chars)
	words := strings.Fields(text)

	words, err = processor.Process(words)
	if err != nil {
		fmt.Println("Error")
		return
	}

	words = processor.Clean(words)

	result := strings.Join(words, " ")
	result = processor.FixPunctuation(result)
	result = processor.FixQuotes(result)

	os.WriteFile(os.Args[2], []byte(result), 0644)
}

package inpt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

// Line reads a line from stdin and trims the whitespace around it.
func Line() (string, error) {
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("could not read line: %w", err)
	}
	result := scanner.Text()
	result = strings.TrimSpace(result)
	return result, nil
}

// Confirm reads input from the user and returns true if it is "y" or false if
// it is anything else.
func Confirm() (bool, error) {
	input, err := Line()
	if err != nil {
		return false, fmt.Errorf("could not get user confirmation: %w", err)
	}
	input = strings.ToLower(input)
	if input != "y" {
		return false, nil
	}
	return true, nil
}

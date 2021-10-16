package inpt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var scanner *Scanner

func init() {
	scanner = NewScanner(os.Stdin)
}

// Normalize trims the leading and trailing whitespace of a string and makes all
// of its characters lowercase.
func Normalize(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	return str
}

// Scanner is a wrapper around bufio.Scanner with convenient methods for
// collecting user input and cleaning it up.
type Scanner struct {
	*bufio.Scanner
}

// NewScanner returns a Scanner of the supplied reader.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{bufio.NewScanner(r)}
}

// Line reads a line from the scanner and trims the whitespace around it.
func (s *Scanner) Line() (string, error) {
	s.Scan()
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("could not read line: %w", err)
	}
	result := scanner.Text()
	result = strings.TrimSpace(result)
	return result, nil
}

// Confirm reads input from the scanner and returns true if it is "y" or false if
// it is anything else.
func (s *Scanner) Confirm() (bool, error) {
	input, err := s.Line()
	if err != nil {
		return false, fmt.Errorf("could not get user confirmation: %w", err)
	}
	input = strings.ToLower(input)
	if input != "y" {
		return false, nil
	}
	return true, nil
}

// Line reads a line from stdin and trims the whitespace around it.
func Line() (string, error) {
	return scanner.Line()
}

// Confirm reads input from the user and returns true if it is "y" or false if
// it is anything else.
func Confirm() (bool, error) {
	return scanner.Confirm()
}

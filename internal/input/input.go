package input

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
		return "", fmt.Errorf("input: error reading line from user: %w", err)
	}
	result := scanner.Text()
	result = strings.TrimSpace(result)
	return result, nil
}

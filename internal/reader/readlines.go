package reader

import (
	"bufio"
	"os"
)

// ReadLines reads a file and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	var lines []string
	inFile, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

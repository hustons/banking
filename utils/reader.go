package utils

import (
	"bufio"
	"log"
	"os"
)

func ReadInputFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Could not open file ", filename, ": ", err)
		return nil
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Failed while reading file ", filename, ": ", err)
	}

	return lines
}

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	printInfo("Expected args are file path, match text, replace text")

	if len(os.Args) < 4 {
		printError("Invalid number of arguments passed")
		return
	}

	filePath := os.Args[1]
	match := os.Args[2]
	replace := os.Args[3]

	if !strings.Contains(replace, "%s") {
		printError("Replace string must include '%s' for GUID placement")
		return
	}

	contents, err := os.ReadFile(filePath)

	if err != nil && errors.Is(err, os.ErrNotExist) {
		printError("Invalid file path")
		return
	}

	if err != nil {
		printFormattedError("Error reading file: %s", err)
		return
	}

	if len(contents) == 0 {
		printError("File has no contents")
		return
	}

	parts := strings.Split(string(contents), match)

	if len(parts) < 2 {
		printError("File has no matching content to replace")
		return
	}

	var sb strings.Builder

	sb.WriteString(parts[0])

	for _, part := range parts[1:] {
		sb.WriteString(fmt.Sprintf(replace, uuid.New()))
		sb.WriteString(part)
	}

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		printFormattedError("Error getting file permissions: %s", err)
		return
	}

	err = os.WriteFile(filePath, []byte(sb.String()), fileInfo.Mode())

	if err != nil {
		printFormattedError("Error writing to file: %s", err)
		return
	}

	printFormattedInfo("Successfully replaced matches with %d Guids", len(parts)-1)
}

func printFormattedError(line string, values ...any) {
	printError(fmt.Sprintf(line, values...))
}

func printFormattedInfo(line string, values ...any) {
	printInfo(fmt.Sprintf(line, values...))
}

func printError(line string) {
	fmt.Println("ERROR: " + line)
}

func printInfo(line string) {
	fmt.Println("INFO: " + line)
}

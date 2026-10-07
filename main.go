package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <entry_file> <exit_file>")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Reading error :", err)
		return
	}

	transformedText := ProcessText(string(content))

	err = os.WriteFile(os.Args[2], []byte(transformedText), 0644)
	if err != nil {
		fmt.Println("Writing error :", err)
		return
	}

	fmt.Println("File treated")
}

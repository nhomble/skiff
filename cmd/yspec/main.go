package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"yspec/pkg/diff"
	"yspec/pkg/k8s"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <before.yaml> <after.yaml>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Use '-' for stdin\n")
		os.Exit(1)
	}

	beforePath := os.Args[1]
	afterPath := os.Args[2]

	beforeReader, err := openInput(beforePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening before stream: %v\n", err)
		os.Exit(1)
	}
	defer closeInput(beforeReader, beforePath)

	afterReader, err := openInput(afterPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening after stream: %v\n", err)
		os.Exit(1)
	}
	defer closeInput(afterReader, afterPath)

	beforeObjects, err := k8s.ParseYAMLStream(beforeReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing before stream: %v\n", err)
		os.Exit(1)
	}

	afterObjects, err := k8s.ParseYAMLStream(afterReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing after stream: %v\n", err)
		os.Exit(1)
	}

	result, err := diff.Generate(beforeObjects, afterObjects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating diff: %v\n", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding output: %v\n", err)
		os.Exit(1)
	}
}

// openInput opens a file or returns stdin if path is "-"
func openInput(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

// closeInput closes a file if it's not stdin
func closeInput(reader io.Reader, path string) {
	if path != "-" {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}
}
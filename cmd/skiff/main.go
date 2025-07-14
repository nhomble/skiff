package main

import (
	"encoding/json"
	"fmt"
	"os"

	"skiff/pkg/diff"
	"skiff/pkg/k8s"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <before.yaml> <after.yaml>\n", os.Args[0])
		os.Exit(1)
	}

	beforePath := os.Args[1]
	afterPath := os.Args[2]

	beforeFile, err := os.Open(beforePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", beforePath, err)
		os.Exit(1)
	}
	defer beforeFile.Close() // nolint:errcheck

	afterFile, err := os.Open(afterPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", afterPath, err)
		os.Exit(1)
	}
	defer afterFile.Close() // nolint:errcheck

	beforeObjects, err := k8s.ParseYAMLStream(beforeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", beforePath, err)
		os.Exit(1)
	}

	afterObjects, err := k8s.ParseYAMLStream(afterFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", afterPath, err)
		os.Exit(1)
	}

	result, err := diff.GenerateTerraformStyle(beforeObjects, afterObjects)
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

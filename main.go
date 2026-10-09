package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "go-reloaded:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: go run . <input.txt> <output.txt>")
	}
	input, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	output, err := Format(string(input))
	if err != nil {
		return err
	}
	if err := os.WriteFile(args[1], []byte(output), 0644); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

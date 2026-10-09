package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "input.txt"), filepath.Join(dir, "output.txt")
	if err := os.WriteFile(input, []byte("a amazing day (up, 2) !"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input, output}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(output)
	if err != nil || string(got) != "an AMAZING DAY!" {
		t.Fatalf("output = %q, error = %v", got, err)
	}
	if err := run([]string{input, input}); err != nil {
		t.Fatal(err)
	}
	got, err = os.ReadFile(input)
	if err != nil || string(got) != "an AMAZING DAY!" {
		t.Fatalf("in-place output = %q, error = %v", got, err)
	}
}

func TestRunErrors(t *testing.T) {
	dir := t.TempDir()
	input, output := filepath.Join(dir, "input.txt"), filepath.Join(dir, "output.txt")
	for _, args := range [][]string{nil, {input}, {input, output, "extra"}, {input, output}} {
		if err := run(args); err == nil {
			t.Errorf("run(%q) should fail", args)
		}
	}
	if err := os.WriteFile(input, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input, dir}); err == nil {
		t.Error("writing to a directory should fail")
	}
	if err := os.WriteFile(input, []byte("wrong (hex)"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("keep me"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{input, output}); err == nil {
		t.Error("invalid input should fail")
	}
	got, err := os.ReadFile(output)
	if err != nil || string(got) != "keep me" {
		t.Fatalf("existing output changed on format failure: %q, %v", got, err)
	}
}

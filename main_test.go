package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGolden(t *testing.T) {
	// Build the plugin
	buildCmd := exec.Command("go", "build", "-o", "protoc-gen-elm", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build protoc-gen-elm: %v", err)
	}

	// Create a temporary directory for the output
	tempDir, err := os.MkdirTemp("", "protoc-elm-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run protoc
	protocCmd := exec.Command("protoc",
		`--plugin=./protoc-gen-elm`,
		`--elm_out=`+tempDir,
		`-I=testdata`,
		`testdata/test.proto`,
	)
	protocCmd.Stdout = os.Stdout
	protocCmd.Stderr = os.Stderr
	if err := protocCmd.Run(); err != nil {
		t.Fatalf("protoc failed: %v", err)
	}

	// Compare the generated file with the golden file
	generatedFile := filepath.Join(tempDir, "Test/Test.elm")
	goldenFile := "testdata/Test.elm.golden"

	generatedBytes, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	goldenBytes, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}

	if !bytes.Equal(generatedBytes, goldenBytes) {
		t.Errorf(`Generated file does not match golden file.\nGenerated:\n%s\nGolden:\n%s`, string(generatedBytes), string(goldenBytes))
	}
}
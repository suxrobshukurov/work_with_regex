package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetFilePaths(t *testing.T) {
	// буфер для подмены ввода
	input := "input.txt\noutput.txt\n"
	expectedInputPath := "input.txt"
	expectedOutputPath := "output.txt"

	// подмена ввода
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // востанавливаем после теста

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdin = r

	go func() {
		w.Write([]byte(input))
		w.Close()
	}()

	inputPath, outputPath, err := getFilePaths()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Проверяем результат
	if inputPath != expectedInputPath {
		t.Errorf("Expected input path %q, got %q", expectedInputPath, inputPath)
	}
	if outputPath != expectedOutputPath {
		t.Errorf("Expected output path %q, got %q", expectedOutputPath, outputPath)
	}
}

func TestReadFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	content := "10+20=?\n"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	readContent, err := readFile(tempFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if readContent != content {
		t.Errorf("Expected %q, got %q", content, readContent)
	}
}

func TestProcessContent(t *testing.T) {
	content := "10+20=?\n15-5=?\n6*7=?\n42/6=?\n100/0=?\n"
	expected := []string{"10+20=30", "15-5=10", "6*7=42", "42/6=7"}

	result, err := processContent(content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPerformOperation(t *testing.T) {
	tests := []struct {
		a, b int
		op   string
		want string
	}{
		{10, 20, "+", "30"},
		{15, 5, "-", "10"},
		{6, 7, "*", "42"},
		{42, 6, "/", "7"},
		{100, 0, "/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			got := performOperation(tt.a, tt.b, tt.op)
			if got != tt.want {
				t.Errorf("performOperation(%d, %d, %q) = %v, want %v", tt.a, tt.b, tt.op, got, tt.want)
			}
		})
	}
}

func TestWriteResults(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test_output_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	results := []string{"10+20=30", "15-5=10", "6*7=42", "42/6=7"}

	err = writeResults(tempFile.Name(), results)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	fileContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := "10+20=30\n15-5=10\n6*7=42\n42/6=7\n"
	if string(fileContent) != expectedContent {
		t.Errorf("Expected %q, got %q", expectedContent, string(fileContent))
	}
}


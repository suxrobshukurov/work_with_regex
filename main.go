package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)


func getFilePaths() (string, string, error) {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите путь к файлу для входных данных: ")
	inputPath, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	inputPath = strings.TrimSpace(inputPath)

	fmt.Print("Введите путь к файлу для выходных данных: ")
	outputPath, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	outputPath = strings.TrimSpace(outputPath)

	if outputPath == "" {
		outputPath = "output.txt"
	}

	return inputPath, outputPath, nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func processContent(content string) ([]string, error) {
	var re = regexp.MustCompile(`([0-9]+)([+\-*/]{1})([0-9]+)([=]{1})([?]{1})`)
	submatch := re.FindAllStringSubmatch(content, -1)
	results := make([]string, 0, len(submatch))

	for _, s := range submatch {
		a, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}

		b, err := strconv.Atoi(s[3])
		if err != nil {
			return nil, err
		}
		result := performOperation(a, b, s[2])
		if result != "" {
			results = append(results, fmt.Sprintf("%v%v%v%v%v", a, s[2], b, s[4], result))
		}
	}

	return results, nil
}

func performOperation(a, b int, op string) string {
	switch op {
	case "+":
		return strconv.Itoa(a + b)
	case "-":
		return strconv.Itoa(a - b)
	case "*":
		return strconv.Itoa(a * b)
	case "/":
		if b != 0 {
			return strconv.Itoa(a / b)
		} else {
			return "Деление на ноль нельзя!"	
		}
	}
}

func writeResults(path string, results []string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)

	for _, r := range results {
		_, err := w.WriteString(r + "\n")
		if err != nil {
			return err
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func main() {
	inputPath, outputPath, err := getFilePaths()

	if err != nil {
		log.Fatal(err)
	}

	content, err := readFile(inputPath)

	if err != nil {
		log.Fatal(err)
	}

	result, err := processContent(content)

	if err != nil {
		log.Fatal(err)
	}

	err = writeResults(outputPath, result)

	if err != nil {
		log.Fatal(err)
	}
}

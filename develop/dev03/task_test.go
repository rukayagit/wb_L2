package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// Утилита для сравнения двух срезов строк
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Функция сортировки для тестирования
func sortStringsForTest(input []string, opts SortOptions) []string {
	return sortStrings(input, opts)
}

func TestSortStrings(t *testing.T) {
	tests := []struct {
		input    []string
		opts     SortOptions
		expected []string
	}{
		{
			input:    []string{"banana", "apple", "cherry"},
			opts:     SortOptions{},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			input:    []string{"2 apple", "1 banana", "3 cherry"},
			opts:     SortOptions{column: 1, numeric: true},
			expected: []string{"1 banana", "2 apple", "3 cherry"},
		},
		{
			input:    []string{"banana", "apple", "banana"},
			opts:     SortOptions{unique: true},
			expected: []string{"apple", "banana"},
		},
		{
			input:    []string{"apple", "banana", "cherry"},
			opts:     SortOptions{reverse: true},
			expected: []string{"cherry", "banana", "apple"},
		},
		{
			input:    []string{"  banana", "apple ", " cherry "},
			opts:     SortOptions{ignoreTrailingSpaces: true},
			expected: []string{"apple", "banana", "cherry"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := sortStringsForTest(test.input, test.opts)
			if !equal(result, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

// Функция для тестирования работы с файлами
func TestFileSort(t *testing.T) {
	// Создание временного входного файла
	inputFile, err := ioutil.TempFile("", "input.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(inputFile.Name())

	// Создание временного выходного файла
	outputFile, err := ioutil.TempFile("", "output.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile.Name())

	// Запись тестовых данных в входной файл
	inputData := []string{"banana", "apple", "cherry"}
	for _, line := range inputData {
		_, err := inputFile.WriteString(line + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	inputFile.Close() // Закрываем файл для завершения записи

	// Запуск функции сортировки с использованием внешних параметров
	opts := SortOptions{}
	err = runSort(inputFile.Name(), outputFile.Name(), opts)
	if err != nil {
		t.Fatal(err)
	}

	// Чтение данных из выходного файла
	outputData, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := "apple\nbanana\ncherry\n"
	if string(outputData) != expectedOutput {
		t.Errorf("expected %s, got %s", expectedOutput, string(outputData))
	}
}

// Утилита для запуска программы с заданными параметрами
func runSort(inputFileName, outputFileName string, opts SortOptions) error {
	args := []string{
		"-input", inputFileName,
		"-output", outputFileName,
	}
	if opts.column > 0 {
		args = append(args, "-k", fmt.Sprintf("%d", opts.column))
	}
	if opts.numeric {
		args = append(args, "-n")
	}
	if opts.reverse {
		args = append(args, "-r")
	}
	if opts.unique {
		args = append(args, "-u")
	}
	if opts.month {
		args = append(args, "-M")
	}
	if opts.ignoreTrailingSpaces {
		args = append(args, "-b")
	}
	if opts.checkSorted {
		args = append(args, "-c")
	}
	if opts.suffix {
		args = append(args, "-h")
	}

	// Эмулируем вызов main с использованием аргументов
	os.Args = append([]string{"sortutil"}, args...)
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
	main()
	return nil
}

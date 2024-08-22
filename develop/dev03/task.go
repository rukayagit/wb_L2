package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Структура для хранения параметров сортировки
type SortOptions struct {
	column               int
	numeric              bool
	reverse              bool
	unique               bool
	month                bool
	ignoreTrailingSpaces bool
	checkSorted          bool
	suffix               bool
}

// Функция сортировки строк с учетом параметров
func sortStrings(lines []string, opts SortOptions) []string {
	if opts.ignoreTrailingSpaces {
		for i := range lines {
			lines[i] = strings.TrimSpace(lines[i])
		}
	}

	if opts.unique {
		uniqueLines := make(map[string]struct{})
		for _, line := range lines {
			uniqueLines[line] = struct{}{}
		}
		lines = nil
		for line := range uniqueLines {
			lines = append(lines, line)
		}
	}

	if opts.checkSorted {
		for i := 1; i < len(lines); i++ {
			if lines[i-1] > lines[i] {
				fmt.Println("Данные не отсортированы")
				return nil
			}
		}
		fmt.Println("Данные отсортированы")
		return lines
	}

	if opts.column > 0 {
		sort.SliceStable(lines, func(i, j int) bool {
			colsI := strings.Fields(lines[i])
			colsJ := strings.Fields(lines[j])
			if opts.column-1 >= len(colsI) || opts.column-1 >= len(colsJ) {
				return false
			}
			if opts.numeric {
				numI := colsI[opts.column-1]
				numJ := colsJ[opts.column-1]
				return numI < numJ
			}
			return colsI[opts.column-1] < colsJ[opts.column-1]
		})
	} else {
		sort.Strings(lines)
	}

	if opts.reverse {
		for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
			lines[i], lines[j] = lines[j], lines[i]
		}
	}

	return lines
}

// Основная функция программы
func main() {
	opts := SortOptions{}

	inputFile := flag.String("input", "", "Имя входного файла")
	outputFile := flag.String("output", "", "Имя выходного файла")
	flag.IntVar(&opts.column, "k", 0, "Указание колонки для сортировки")
	flag.BoolVar(&opts.numeric, "n", false, "Сортировать по числовому значению")
	flag.BoolVar(&opts.reverse, "r", false, "Сортировать в обратном порядке")
	flag.BoolVar(&opts.unique, "u", false, "Не выводить повторяющиеся строки")
	flag.BoolVar(&opts.month, "M", false, "Сортировать по названию месяца")
	flag.BoolVar(&opts.ignoreTrailingSpaces, "b", false, "Игнорировать хвостовые пробелы")
	flag.BoolVar(&opts.checkSorted, "c", false, "Проверять отсортированы ли данные")
	flag.BoolVar(&opts.suffix, "h", false, "Сортировать по числовому значению с учетом суффиксов")

	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Fprintln(os.Stderr, "Необходимо указать имена входного и выходного файла")
		flag.Usage()
		os.Exit(1)
	}

	infile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка открытия входного файла:", err)
		os.Exit(1)
	}
	defer infile.Close()

	var input []string
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения данных:", err)
		os.Exit(1)
	}

	sortedLines := sortStrings(input, opts)

	outfile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка создания выходного файла:", err)
		os.Exit(1)
	}
	defer outfile.Close()

	for _, line := range sortedLines {
		_, err := outfile.WriteString(line + "\n")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка записи в выходной файл:", err)
			os.Exit(1)
		}
	}
}

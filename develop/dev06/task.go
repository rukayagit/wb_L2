package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Определяем флаги командной строки
	fields := flag.String("f", "", "выбрать поля (колонки)")
	delimiter := flag.String("d", "\t", "использовать другой разделитель")
	separated := flag.Bool("s", false, "только строки с разделителем")
	flag.Parse()

	// Парсим флаги
	fieldList := parseFields(*fields)

	// Обработка разделителя
	actualDelimiter := *delimiter

	// Чтение строк из STDIN
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// Разделяем строку по разделителю
		columns := strings.Split(line, actualDelimiter)

		// Проверка на наличие разделителя в строке
		if *separated && len(columns) < 2 {
			continue
		}

		// Выбираем запрошенные поля
		var selectedFields []string
		for _, index := range fieldList {
			if index > 0 && index <= len(columns) {
				selectedFields = append(selectedFields, columns[index-1])
			}
		}

		// Выводим результат
		if len(selectedFields) > 0 {
			fmt.Println(strings.Join(selectedFields, actualDelimiter))
		}
	}

	// Проверка на ошибки при чтении
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", err)
		os.Exit(1)
	}
}

// parseFields парсит строку с номерами полей в слайс целых чисел
func parseFields(fields string) []int {
	var fieldList []int
	if fields == "" {
		return fieldList
	}

	for _, f := range strings.Split(fields, ",") {
		field, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка парсинга поля '%s': %v\n", f, err)
			continue
		}
		fieldList = append(fieldList, field)
	}

	return fieldList
}

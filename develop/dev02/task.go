package main

import (
	"fmt"
	"unicode"
)

// extract распаковывает строку, поддерживая escape-последовательности.
func extract(str string) (string, error) {
	var res []rune
	letters := []rune(str)
	length := len(letters)

	for i := 0; i < length; i++ {
		if letters[i] == '\\' {
			// Если '\' последний символ в строке, это ошибка
			if i+1 >= length {
				return "", fmt.Errorf("Некорректная строка")
			}
			next := letters[i+1]
			// Если следующий символ после '\' снова '\', это экранирование '\'
			if next == '\\' || unicode.IsDigit(next) {
				res = append(res, next)
				i++ // Пропустить следующий символ
			} else {
				return "", fmt.Errorf("Некорректная строка")
			}
		} else if unicode.IsDigit(letters[i]) {
			// Если цифра идет без экранирования или сразу после другой цифры, это ошибка
			if i == 0 || unicode.IsDigit(letters[i-1]) {
				return "", fmt.Errorf("Некорректная строка")
			}
			count := int(letters[i] - '0')
			for j := 0; j < count-1; j++ {
				res = append(res, letters[i-1])
			}
		} else {
			res = append(res, letters[i])
		}
	}

	return string(res), nil
}

func main() {
	tests := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		`qwe\4\5`,
		`qwe\45`,
		`qwe\\5`,
	}

	for _, test := range tests {
		unpacked, err := extract(test)
		if err != nil {
			fmt.Printf("Ошибка распаковки строки: %q: %v\n", test, err)
		} else {
			fmt.Printf("Исходная строка: : %q, Распакованная строка: %q\n", test, unpacked)
		}
	}
}

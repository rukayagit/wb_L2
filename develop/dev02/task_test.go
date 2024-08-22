package main

import (
	"testing"
)

// Тестируем функцию extract с поддержкой escape-последовательностей
func TestExtract(t *testing.T) {
	tests := []struct {
		input    string // Входная строка для распаковки
		expected string // Ожидаемый результат
		hasError bool   // Ожидается ли ошибка
	}{
		{"a4bc2d5e", "aaaabccddddde", false}, // Тест с корректной строкой
		{"abcd", "abcd", false},              // Тест с корректной строкой без цифр
		{"45", "", true},                     // Тест с некорректной строкой (начинается с цифры)
		{"", "", false},                      // Тест с пустой строкой
		{`qwe\4\5`, "qwe45", false},          // Тест с экранированными цифрами
		{`qwe\45`, "qwe44444", false},        // Тест с экранированной цифрой и числом
		{`qwe\\5`, "qwe\\\\\\", false},       // Тест с экранированной обратной косой чертой и числом
		{`\\`, `\`, false},                   // Тест с одной экранированной обратной косой чертой
		{`\`, "", true},                      // Тест с одиночной обратной косой чертой (ошибка)
		{`a\`, "", true},                     // Тест с обратной косой чертой в конце строки (ошибка)
		{`\a`, "", true},                     // Тест с неправильной экранированной последовательностью
	}

	// Пробегаемся по каждому тесту
	for _, test := range tests {
		result, err := extract(test.input) // Распаковываем строку
		if (err != nil) != test.hasError { // Проверяем, совпадает ли наличие ошибки с ожиданием
			t.Errorf("extract(%q) returned error %v, expected error: %v", test.input, err, test.hasError)
		}
		if result != test.expected { // Проверяем, совпадает ли результат с ожидаемым
			t.Errorf("extract(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

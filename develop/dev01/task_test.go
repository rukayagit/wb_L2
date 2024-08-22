package main

import (
	"testing"
)

// Эта функция тестирует нашу программу с различными аргументами и проверяет, правильно ли она работает.
func TestCLI(t *testing.T) {
	tests := []struct {
		name    string   // Имя теста
		args    []string // Аргументы командной строки для теста
		wantErr bool     // Ожидаем ли ошибку
	}{
		{"default host", []string{}, false},                                    // Тест с настройками по умолчанию (не должно быть ошибки).
		{"custom host", []string{"-host", "0.beevik-ntp.pool.ntp.org"}, false}, // Тест с указанным хостом (не должно быть ошибки).
		{"invalid host", []string{"-host", "invalid.host"}, true},              // Тест с неправильным хостом (должна быть ошибка).
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CLI(tt.args); (got != 0) != tt.wantErr {
				t.Errorf("CLI() = %v, wantErr %v", got, tt.wantErr) // Если результат не совпадает с ожидаемым, выводим ошибку.
			}
		})
	}
}

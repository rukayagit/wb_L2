package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// Эта функция запускает нашу программу и возвращает код завершения
// 0 - успешное завершение, 1 - ошибка выполнения, 2 - ошибка в аргументах
func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args) // Разбираем аргументы командной строки
	if err != nil {
		return 2 // Если в аргументах есть ошибка, возвращаем 2
	}

	if err := app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err) // Печатаем ошибку выполнения в STDERR
		return 1                                           // Если есть ошибка выполнения, возвращаем 1
	}
	return 0 // Если нет ошибок, возвращаем 0
}

// Это структура для хранения параметров приложения.
// Пока здесь только один параметр — хост (адрес) NTP-сервера.
type appEnv struct {
	host string
}

// Эта функция разбирает аргументы командной строки и сохраняет их в структуре appEnv.
func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("ntp", flag.ContinueOnError)
	fl.StringVar(&app.host, "host", "0.beevik-ntp.pool.ntp.org", "ntp host")

	if err := fl.Parse(args); err != nil {
		return err // Если есть ошибка в аргументах, возвращаем её.
	}

	return nil // Если аргументы успешно разобраны, возвращаем nil
}

// Эта функция получает текущее время с указанного NTP-сервера и выводит его.
func (app *appEnv) run() error {
	t, err := ntp.Time(app.host) // Запрашиваем текущее время у NTP-сервера.
	if err != nil {
		return err // Если есть ошибка при запросе времени, возвращаем её.
	}

	fmt.Println(t.UTC().Format(time.UnixDate)) // Выводим текущее время в формате UnixDate
	return nil                                 // Если нет ошибок, возвращаем nil
}

// Это точка входа в программу. Она вызывает функцию CLI с аргументами командной строки и завершает программу с соответствующим кодом.
func main() {
	os.Exit(CLI(os.Args[1:]))
}

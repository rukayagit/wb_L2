package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Объявляем глобальные переменные для URL и директории сохранения
var (
	startURL  string
	outputDir string
)

// Функция init инициализирует флаги командной строки
func init() {
	// Определяем флаг для URL, который нужно загрузить
	flag.StringVar(&startURL, "url", "", "URL для загрузки")
	// Определяем флаг для директории, в которую будут сохранены файлы
	flag.StringVar(&outputDir, "output", ".", "Директория для сохранения загруженных файлов")
}

// Функция main является точкой входа в программу
func main() {
	// Обрабатываем (парсим) флаги командной строки
	flag.Parse()

	// Если URL не указан, выводим сообщение и завершаем работу
	if startURL == "" {
		fmt.Println("Пожалуйста, укажите URL с помощью флага -url")
		return
	}

	// Создаем директорию для сохранения файлов, если она не существует
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return
	}

	// Пытаемся загрузить страницу по указанному URL и сохранить её в указанную директорию
	if err := downloadPage(startURL, outputDir); err != nil {
		fmt.Printf("Ошибка загрузки страницы: %v\n", err)
	}
}

// downloadPage загружает страницу по указанному URL и сохраняет её в указанной директории
func downloadPage(pageURL, baseDir string) error {
	// Выполняем HTTP GET-запрос на указанный URL
	resp, err := http.Get(pageURL)
	if err != nil {
		// Возвращаем ошибку, если запрос не удался
		return fmt.Errorf("не удалось получить URL %s: %w", pageURL, err)
	}
	// Закрываем ответ (ресурс) после завершения работы с ним
	defer resp.Body.Close()

	// Определяем путь для сохранения файла на основе URL и базовой директории
	pagePath := getFilePath(pageURL, baseDir)

	// Создаем все необходимые директории, если они не существуют
	if err := os.MkdirAll(filepath.Dir(pagePath), 0755); err != nil {
		// Возвращаем ошибку, если создание директории не удалось
		return fmt.Errorf("не удалось создать директории: %w", err)
	}

	// Создаем файл для сохранения содержимого страницы
	file, err := os.Create(pagePath)
	if err != nil {
		// Возвращаем ошибку, если файл не удалось создать
		return fmt.Errorf("не удалось создать файл %s: %w", pagePath, err)
	}
	// Закрываем файл после завершения работы с ним
	defer file.Close()

	// Копируем содержимое HTTP-ответа (страницы) в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// Возвращаем ошибку, если не удалось записать содержимое страницы
		return fmt.Errorf("не удалось записать HTML: %w", err)
	}

	// Выводим сообщение об успешной загрузке страницы
	fmt.Printf("Загружено %s в %s\n", pageURL, pagePath)
	return nil
}

// getFilePath преобразует URL в путь для сохранения файла на диске
func getFilePath(pageURL, baseDir string) string {
	// Парсим URL, чтобы извлечь путь
	parsedURL, _ := url.Parse(pageURL)
	path := parsedURL.Path

	// Если путь пустой или равен "/", сохраняем как "index.html"
	if path == "" || path == "/" {
		path = "index.html"
	} else {
		// Удаляем начальный слэш, если он есть
		path = strings.TrimPrefix(path, "/")
	}

	// Объединяем базовую директорию с полученным путем, чтобы получить полный путь к файлу
	return filepath.Join(baseDir, path)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Event представляет собой событие в календаре
type Event struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
}

// Глобальная переменная для хранения событий
var events = make(map[int]Event)

// Вспомогательная функция для парсинга строки даты в формат time.Time
func parseDate(dateStr string) (time.Time, error) {
	// Формат даты - "2006-01-02"
	return time.Parse("2006-01-02", dateStr)
}

// Вспомогательная функция для отправки JSON-ответов
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Обработчик для создания события
func createEventHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение id из запроса
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	// Извлечение title из запроса
	title := r.FormValue("title")

	// Извлечение user_id из запроса
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	// Извлечение и парсинг даты из запроса
	date, err := parseDate(r.FormValue("date"))
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	// Сохранение события
	events[id] = Event{ID: id, Title: title, UserID: userID, Date: date}

	// Логирование события после создания
	log.Printf("Event created: %+v", events[id])

	// Отправка ответа
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "event created"})
}

// Обработчик для обновления события
func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение id из запроса
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	// Проверка существования события
	event, exists := events[id]
	if !exists {
		http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
		return
	}

	// Извлечение title из запроса
	title := r.FormValue("title")

	// Извлечение user_id из запроса
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid user_id"}`, http.StatusBadRequest)
		return
	}

	// Извлечение и парсинг даты из запроса
	date, err := parseDate(r.FormValue("date"))
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	// Обновление события
	event.Title = title
	event.UserID = userID
	event.Date = date
	events[id] = event

	// Логирование события после обновления
	log.Printf("Event updated: %+v", events[id])

	// Отправка ответа
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "event updated"})
}

// Обработчик для удаления события
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение id из запроса
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, `{"error": "Invalid id"}`, http.StatusBadRequest)
		return
	}

	// Проверка существования события
	if _, exists := events[id]; !exists {
		http.Error(w, `{"error": "Event not found"}`, http.StatusNotFound)
		return
	}

	// Удаление события
	delete(events, id)

	// Логирование после удаления события
	log.Printf("Event with ID %d deleted", id)

	// Отправка ответа
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "event deleted"})
}

// Обработчик для получения событий за день
func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение даты из запроса
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	var result []Event
	for _, event := range events {
		if event.Date.Equal(date) {
			result = append(result, event)
		}
	}

	// Логирование результатов поиска
	log.Printf("Events for day %s: %+v", dateStr, result)

	// Отправка ответа
	if len(result) == 0 {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

// Обработчик для получения событий за неделю
func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение даты из запроса
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	// Определение начала и конца недели
	startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var result []Event
	for _, event := range events {
		if event.Date.After(startOfWeek) && event.Date.Before(endOfWeek) {
			result = append(result, event)
		}
	}

	// Логирование результатов поиска
	log.Printf("Events for week starting %s: %+v", startOfWeek.Format("2006-01-02"), result)

	// Отправка ответа
	if len(result) == 0 {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

// Обработчик для получения событий за месяц
func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение даты из запроса
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date"}`, http.StatusBadRequest)
		return
	}

	// Определение начала и конца месяца
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var result []Event
	for _, event := range events {
		if (event.Date.After(startOfMonth) || event.Date.Equal(startOfMonth)) && (event.Date.Before(endOfMonth) || event.Date.Equal(endOfMonth)) {
			result = append(result, event)
		}
	}

	// Логирование результатов поиска
	log.Printf("Events for month %s: %+v", dateStr, result)

	// Отправка ответа
	if len(result) == 0 {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

// Middleware для логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func main() {
	// Создание нового мультиплексора маршрутов
	mux := http.NewServeMux()
	// Регистрация обработчиков
	mux.HandleFunc("/create_event", createEventHandler)
	mux.HandleFunc("/update_event", updateEventHandler)
	mux.HandleFunc("/delete_event", deleteEventHandler)
	mux.HandleFunc("/events_for_day", eventsForDayHandler)
	mux.HandleFunc("/events_for_week", eventsForWeekHandler)
	mux.HandleFunc("/events_for_month", eventsForMonthHandler)

	// Запуск сервера
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", loggingMiddleware(mux))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

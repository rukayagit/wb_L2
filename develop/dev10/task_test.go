package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// startTestServer запускает тестовый TCP сервер на заданном порту.
// Параметр t используется для логирования и обработки ошибок в тестах.
func startTestServer(t *testing.T, port string) net.Listener {
	// Открываем TCP-соединение на указанном порту.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// Если произошла ошибка при попытке открыть соединение, завершаем тест с ошибкой.
		t.Fatalf("Error starting test server: %v", err)
	}

	// Логируем сообщение о том, что сервер успешно запущен и слушает на порту.
	t.Logf("Test server listening on port %s", port)

	// Запускаем горутину для обработки входящих соединений.
	go func() {
		for {
			// Ожидаем входящего соединения.
			conn, err := listener.Accept()
			if err != nil {
				// Логируем ошибку при принятии соединения и завершаем обработку.
				t.Logf("Error accepting connection: %v", err)
				return
			}
			// Обрабатываем подключенного клиента в новой горутине.
			go handleTestClientConnection(t, conn)
		}
	}()

	// Возвращаем объект listener для управления сервером (например, для закрытия).
	return listener
}

// handleTestClientConnection обрабатывает входящие сообщения от клиента и отправляет обратно их же, но в верхнем регистре.
func handleTestClientConnection(t *testing.T, conn net.Conn) {
	defer conn.Close()              // Закрываем соединение по завершении функции.
	reader := bufio.NewReader(conn) // Создаем буферизированный ридер для чтения данных от клиента.

	for {
		// Читаем строку от клиента до символа новой строки.
		message, err := reader.ReadString('\n')
		if err != nil {
			// Логируем ошибку при чтении и завершаем обработку соединения.
			t.Logf("Error reading from client: %v", err)
			return
		}

		// Логируем полученное сообщение от клиента.
		t.Logf("Received from client: %s", message)

		// Преобразуем сообщение в верхний регистр и удаляем лишние пробелы.
		newMessage := strings.ToUpper(strings.TrimSpace(message))
		// Отправляем преобразованное сообщение обратно клиенту.
		conn.Write([]byte(newMessage + "\n"))
	}
}

// TestClient тестирует работу клиента, подключающегося к серверу.
func TestClient(t *testing.T) {
	port := "8083"                     // Устанавливаем порт, на котором будет работать сервер.
	server := startTestServer(t, port) // Запускаем тестовый сервер.
	defer server.Close()               // Закрываем сервер по завершении теста.

	time.Sleep(1 * time.Second) // Даем время серверу полностью запуститься.

	address := "localhost:" + port        // Формируем адрес для подключения клиента.
	conn, err := net.Dial("tcp", address) // Подключаемся к серверу.
	if err != nil {
		// Если не удалось подключиться, завершаем тест с ошибкой.
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close() // Закрываем соединение по завершении теста.

	// Сообщения, которые клиент отправит серверу.
	messages := []string{"hello", "how are you?"}
	// Ожидаемые ответы от сервера.
	expectedResponses := []string{"HELLO", "HOW ARE YOU?"}

	for i, msg := range messages {
		// Отправляем сообщение на сервер.
		fmt.Fprintf(conn, msg+"\n")

		// Читаем ответ от сервера.
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// Если произошла ошибка при чтении ответа, завершаем тест с ошибкой.
			t.Fatalf("Error reading from server: %v", err)
		}
		// Удаляем лишние пробелы в ответе.
		response = strings.TrimSpace(response)

		// Проверяем, соответствует ли ответ ожидаемому.
		if response != expectedResponses[i] {
			t.Errorf("Expected response %s but got %s", expectedResponses[i], response)
		}
	}
}

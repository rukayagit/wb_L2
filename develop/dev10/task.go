package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	mode := flag.String("mode", "", "start as server or client (values: server, client)")
	host := flag.String("host", "localhost", "hostname or IP address")
	port := flag.String("port", "8081", "port number")
	timeout := flag.String("timeout", "10s", "timeout for client connection")
	flag.Parse()

	if *mode == "server" {
		startServer(*host, *port)
	} else if *mode == "client" {
		startClient(*host, *port, *timeout)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  To start server: go run main.go --mode=server --port=8081")
		fmt.Println("  To start client: go run main.go --mode=client --host=localhost --port=8081 --timeout=10s")
	}
}

func startServer(host, port string) {
	address := host + ":" + port
	fmt.Println("Launching server on", address)

	// Запуск сервера
	serv, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
		return
	}
	defer serv.Close()

	for {
		fmt.Println("Waiting for connection...")

		// Подключение к каналу
		conn, err := serv.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting: %v\n", err)
			return
		}
		fmt.Println("New connection. Waiting for messages...")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Serving new connection %v\n", conn)

	connReader := bufio.NewReader(conn) // Ридер создается один раз

	for {
		// Чтение сообщения
		message, err := connReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection %v closed.\n", conn)
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
			break
		}
		message = strings.TrimSpace(message) // Удаляем лишние символы

		fmt.Printf("From %v received: %s\n", conn, string(message))

		// Отправка нового сообщения обратно клиенту
		newmessage := strings.ToUpper(message)
		_, err = conn.Write([]byte(newmessage + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to connection: %v\n", err)
			break
		}
	}

	fmt.Printf("Done serving client %v\n", conn)
}

func startClient(host, port, timeout string) {
	exitChan := make(chan os.Signal, 1)
	// При нажатии ctrl + D в канал sigCh будет отправлено сообщение
	signal.Notify(exitChan, syscall.SIGQUIT)
	go func() {
		<-exitChan
		fmt.Println("Press ctrl + D to exit")
		os.Exit(0)
	}()

	// Переводим задержку в нужный формат
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("Error parsing timeout:", err)
		return
	}

	// Формируем строку подключения
	hostPort := host + ":" + port

	// Подключаемся к сокету. Задаем таймаут подключения
	conn, err := net.DialTimeout("tcp", hostPort, timeoutDuration)
	if err != nil {
		fmt.Println("Error connecting to", hostPort, ":", err)
		return
	}
	defer conn.Close()

	// Создаем ридеры
	console := bufio.NewReader(os.Stdin)
	connReader := bufio.NewReader(conn)

	for {
		fmt.Print("Your message: ")

		// Чтение сообщения с консоли
		text, err := console.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading message: %v", err)
			return
		}
		text = strings.TrimSpace(text) // Удаляем лишние символы

		// Отправляем сообщение
		fmt.Fprintf(conn, text+"\n")
		if text == "exit" {
			fmt.Println("Closing connection...")
			return
		}

		// Получаем ответ
		message, err := connReader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response from server: %v", err)
			return
		}

		// Удаляем лишние символы
		message = strings.TrimSpace(message)
		fmt.Printf("From server: %s\n", message)
	}
}

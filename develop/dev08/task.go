package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var ErrNoPath = errors.New("path required")

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Выводим приглашение для ввода команды
		fmt.Print("> ")
		// Считываем строку, введённую пользователем
		input, err := reader.ReadString('\n')
		if err != nil {
			// Если возникла ошибка при чтении, выводим её и продолжаем
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Выполняем команду
		if err = execInput(input); err != nil {
			// Если возникла ошибка при выполнении команды, выводим её
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// Убираем пробельные символы в начале и конце строки
	input = strings.TrimSpace(input)
	// Разделяем строку на аргументы по пробелам
	args := strings.Split(input, " ")

	// Проверяем первую часть команды (основную команду)
	switch args[0] {
	case "cd":
		// Команда для смены директории
		if len(args) < 2 {
			return ErrNoPath
		}
		return os.Chdir(args[1])
	case "pwd":
		// Команда для отображения текущего пути
		if dir, err := os.Getwd(); err == nil {
			fmt.Println(dir)
		} else {
			return err
		}
	case "echo":
		// Команда для вывода аргументов на экран
		fmt.Println(strings.Join(args[1:], " "))
	case "kill":
		// Команда для завершения процесса по PID
		if len(args) < 2 {
			return errors.New("pid required")
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("invalid pid")
		}
		return killProcess(pid)
	case "ps":
		// Команда для отображения списка процессов
		cmd := exec.Command("ps")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	case "exit":
		// Команда для выхода из программы
		os.Exit(0)
	default:
		// Проверка на наличие пайпов в команде
		if strings.Contains(input, "|") {
			return execPipeline(input)
		} else {
			// Выполнение обычной команды
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}

	return nil
}

// Функция для завершения процесса по PID
func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("could not find process with PID %d: %v", pid, err)
	}
	return process.Kill()
}

// Функция для выполнения конвейера (пайпов)
func execPipeline(input string) error {
	// Разделяем команды по символу пайпа
	commands := strings.Split(input, "|")
	var cmds []*exec.Cmd

	// Создаем команды для каждого элемента пайплайна
	for _, cmdStr := range commands {
		args := strings.Fields(strings.TrimSpace(cmdStr))
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	// Связываем stdout одной команды со stdin следующей
	for i := 0; i < len(cmds)-1; i++ {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		cmds[i+1].Stdin = stdout
	}

	// Настраиваем вывод последней команды на stdout и stderr
	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	// Запускаем все команды
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	// Ожидаем завершения всех команд
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

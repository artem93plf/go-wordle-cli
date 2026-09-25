// Package cli реализует консольный интерфейс игры.
package cli

import (
	"bufio"
	"fiveletters/game"
	"fmt"
	"io"
	"strings"
)

// printLine выводит строку в поток вывода.
func printLine(output io.Writer, text string) error {
	_, err := fmt.Fprintln(output, text)
	return err
}

// marksToString преобразует массив результатов проверки буквы в строку.
func marksToString(marks []game.LetterMark) string {
	var builder strings.Builder

	for _, mark := range marks {
		builder.WriteString(string(mark))
	}
	return builder.String()
}

// playGame проигрывает игру с последовательным вводом слов.
func playGame(scanner *bufio.Scanner, output io.Writer, seed int64) error {
	session, err := game.StartGame(game.DefaultDictionary(), 6, seed)
	if err != nil {
		return err
	}

	history := make([]string, 0, session.MaxAttempts)

	for session.Status == game.StatusGameInProgress {
		if err := printLine(output, "Введите слово:"); err != nil {
			return err
		}

		if !scanner.Scan() {
			return scanner.Err()
		}

		result, err := game.ApplyGuess(&session, scanner.Text())
		if err != nil {
			if printErr := printLine(output, err.Error()); printErr != nil {
				return printErr
			}

			continue
		}

		// Добавляем текущую попытку в историю.
		line := fmt.Sprintf("%s %s", marksToString(result.Marks), result.Guess)
		history = append(history, line)

		// Выводим всю историю попыток.
		for _, historyLine := range history {
			if err := printLine(output, historyLine); err != nil {
				return err
			}
		}

		// Показываем оставшееся количество попыток.
		remainedAttempts := session.MaxAttempts - session.UsedAttempts
		if err := printLine(
			output,
			fmt.Sprintf("Осталось попыток: %d", remainedAttempts),
		); err != nil {
			return err
		}

		// Проверяем победу.
		if session.Status == game.StatusGameWin {
			if err := printLine(
				output,
				fmt.Sprintf(
					"Победа! Слово угадано за %d попытки",
					session.UsedAttempts,
				),
			); err != nil {
				return err
			}
		}

		// Проверяем поражение.
		if session.Status == game.StatusGameLose {
			if err := printLine(
				output,
				fmt.Sprintf(
					"Неудача. Загаданное слово: %s",
					session.Answer,
				),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

// Run запускает консольный интерфейс игры.
func Run(input io.Reader, output io.Writer, seed int64) error {
	scanner := bufio.NewScanner(input)

	for {
		if err := printLine(output, "1. Новая игра"); err != nil {
			return err
		}
		if err := printLine(output, "2. Выход"); err != nil {
			return err
		}
		if err := printLine(output, "Выберите опцию:"); err != nil {
			return err
		}

		if !scanner.Scan() {
			return scanner.Err()
		}

		switch scanner.Text() {
		case "1":
			if err := playGame(scanner, output, seed); err != nil {
				return err
			}
		case "2":
			if err := printLine(output, "Выход"); err != nil {
				return err
			}

			return nil
		default:
			if err := printLine(output, "Неверный выбор"); err != nil {
				return err
			}
		}
	}
}

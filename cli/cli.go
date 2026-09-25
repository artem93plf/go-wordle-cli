// Package cli реализует консольный интерфейс игры.
package cli

import (
	"bufio"
	"fiveletters/game"
	"fmt"
	"io"
	"strings"
)

// Colors описывает цвета вывода в консоль.
const (
	colorGreen string = "\033[32m"
	colorRed   string = "\033[31m"
	colorReset string = "\033[0m"
)

// Statistics хранит статистику игры.
type Statistics struct {
	Wins   int
	Losses int
}

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
func playGame(scanner *bufio.Scanner, output io.Writer, seed int64, statistic *Statistics) error {
	category, err := chooseCategory(scanner, output)
	if err != nil {
		return err
	}
	if category == game.CategoryUnknown {
		return nil
	}
	difficulty, err := chooseDifficulty(scanner, output)
	if err != nil {
		return err
	}
	if difficulty == game.DifficultyUnknown {
		return nil
	}
	dictionary := game.DictionaryByOptions(category, difficulty)

	session, err := game.StartGame(dictionary, 6, seed)
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
					colorGreen+"Победа! Слово угадано за %d попытки"+colorReset,
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
					colorRed+"Неудача. Загаданное слово: %s"+colorReset,
					session.Answer,
				),
			); err != nil {
				return err
			}
		}
	}
	if session.Status == game.StatusGameWin {
		statistic.Wins++
	}
	if session.Status == game.StatusGameLose {
		statistic.Losses++
	}
	return nil
}

// chooseCategory предлагает выбрать категорию из списка.
func chooseCategory(scanner *bufio.Scanner, output io.Writer) (game.Category, error) {
	for {
		if err := printLine(output, "Выберите категорию:\n 1. Животные\n 2. Города\n 3. Едa"); err != nil {
			return game.CategoryUnknown, err
		}
		if !scanner.Scan() {
			return game.CategoryUnknown, scanner.Err()
		}

		switch scanner.Text() {
		case "1":
			return game.CategoryAnimal, nil
		case "2":
			return game.CategoryCity, nil
		case "3":
			return game.CategoryFood, nil
		default:
			if err := printLine(output, "Неверный выбор"); err != nil {
				return game.CategoryUnknown, err
			}
		}
	}
}

func chooseDifficulty(scanner *bufio.Scanner, output io.Writer) (game.Difficulty, error) {
	for {
		if err := printLine(output, "Выберите сложность:\n 1. Лёгкая\n 2. Сложная"); err != nil {
			return game.DifficultyUnknown, err
		}
		if !scanner.Scan() {
			return game.DifficultyUnknown, scanner.Err()
		}

		switch scanner.Text() {
		case "1":
			return game.DifficultyEasy, nil
		case "2":
			return game.DifficultyHard, nil
		default:
			if err := printLine(output, "Неверный выбор"); err != nil {
				return game.DifficultyUnknown, err
			}
		}
	}
}

// Run запускает консольный интерфейс игры.
func Run(input io.Reader, output io.Writer, seed int64) error {
	scanner := bufio.NewScanner(input)

	statistic := Statistics{
		Wins:   0,
		Losses: 0,
	}

	for {
		if err := printLine(output, "1. Новая игра"); err != nil {
			return err
		}
		if err := printLine(output, "2. Выход"); err != nil {
			return err
		}
		if err := printLine(output, "3. Статистика"); err != nil {
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
			if err := playGame(scanner, output, seed, &statistic); err != nil {
				return err
			}
		case "2":
			if err := printLine(output, "Выход"); err != nil {
				return err
			}
			return nil
		case "3":
			if err := printLine(output, "Статистика:"); err != nil {
				return err
			}
			if err := printLine(output, fmt.Sprintf("Побед: %d", statistic.Wins)); err != nil {
				return err
			}
			if err := printLine(output, fmt.Sprintf("Поражений: %d", statistic.Losses)); err != nil {
				return err
			}
			if err := printLine(output, fmt.Sprintf("Всего игр: %d", statistic.Wins+statistic.Losses)); err != nil {
				return err
			}
		default:
			if err := printLine(output, "Неверный выбор"); err != nil {
				return err
			}
		}
	}
}

package cli

import (
	"bufio"
	"bytes"
	"fiveletters/game"
	"fmt"
	"strings"
	"testing"
)

func TestRunStatistics(t *testing.T) {
	const win = "1\n3\n1\nсахар\n"
	lose := "1\n3\n1\n" + strings.Repeat("банан\n", 6)
	tests := []struct {
		name   string
		games  string
		wins   int
		losses int
	}{
		{name: "до первой игры"},
		{name: "после победы", games: win, wins: 1},
		{name: "после поражения", games: lose, losses: 1},
		{name: "после нескольких игр", games: win + lose + win, wins: 2, losses: 1},
		{name: "игра после просмотра статистики", games: "3\n" + win, wins: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			// Два просмотра подряд не должны менять счётчики или завершать приложение.
			if err := Run(strings.NewReader(tt.games+"3\n3\n2\n"), &output, 42); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			menu := "1. Новая игра\n2. Выход\n3. Статистика\nВыберите опцию:\n"
			statistics := fmt.Sprintf("Статистика:\nПобед: %d\nПоражений: %d\nВсего игр: %d\n",
				tt.wins, tt.losses, tt.wins+tt.losses)
			want := statistics + menu + statistics + menu + "Выход\n"
			if !strings.HasSuffix(output.String(), want) {
				t.Errorf("Run() = %q, want suffix %q", output.String(), want)
			}
		})
	}
}

func TestRunStatisticsResetBetweenRuns(t *testing.T) {
	var first bytes.Buffer
	if err := Run(strings.NewReader("1\n3\n1\nсахар\n3\n2\n"), &first, 42); err != nil {
		t.Fatalf("first Run() error = %v", err)
	}
	if !strings.Contains(first.String(), "Статистика:\nПобед: 1\nПоражений: 0\nВсего игр: 1\n") {
		t.Fatalf("первый запуск не сохранил победу: %q", first.String())
	}
	var second bytes.Buffer
	if err := Run(strings.NewReader("3\n2\n"), &second, 42); err != nil {
		t.Fatalf("second Run() error = %v", err)
	}
	if !strings.Contains(second.String(), "Статистика:\nПобед: 0\nПоражений: 0\nВсего игр: 0\n") {
		t.Errorf("статистика не сброшена при новом запуске: %q", second.String())
	}
}

func TestRunOptions(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		category   game.Category
		difficulty game.Difficulty
	}{
		{"животные легко", "1\n1\n", game.CategoryAnimal, game.DifficultyEasy},
		{"животные сложно", "1\n2\n", game.CategoryAnimal, game.DifficultyHard},
		{"города легко", "2\n1\n", game.CategoryCity, game.DifficultyEasy},
		{"города сложно", "2\n2\n", game.CategoryCity, game.DifficultyHard},
		{"еда легко", "3\n1\n", game.CategoryFood, game.DifficultyEasy},
		{"еда сложно", "3\n2\n", game.CategoryFood, game.DifficultyHard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := game.StartGame(game.DictionaryByOptions(tt.category, tt.difficulty), 6, 42)
			if err != nil {
				t.Fatalf("StartGame() error = %v", err)
			}
			input := strings.NewReader("1\n" + tt.input + session.Answer + "\n2\n")
			var output bytes.Buffer
			if err := Run(input, &output, 42); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			want := "✅✅✅✅✅ " + session.Answer + "\nОсталось попыток: 5\n\033[32mПобеда! Слово угадано за 1 попытки\033[0m\n"
			if !strings.Contains(output.String(), want) {
				t.Errorf("Run() = %q, want substring %q", output.String(), want)
			}
		})
	}
}

func TestRunInvalidOptions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"неверная категория", "1\n0\n3\n1\nсахар\n2\n"},
		{"неверная сложность", "1\n3\n0\n1\nсахар\n2\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := Run(strings.NewReader(tt.input), &output, 42); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if got := strings.Count(output.String(), "Неверный выбор\n"); got != 1 {
				t.Errorf("число сообщений о неверном выборе = %d, want 1", got)
			}
			if !strings.Contains(output.String(), "\033[32mПобеда! Слово угадано за 1 попытки\033[0m\n") {
				t.Errorf("после исправления выбора игра не завершилась победой: %q", output.String())
			}
		})
	}
}

func TestRunEOFWhileChoosingOptions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"выбор категории", "1\n"},
		{"выбор сложности", "1\n3\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := Run(strings.NewReader(tt.input), &output, 42); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if strings.Contains(output.String(), "Введите слово:") {
				t.Errorf("игра началась без выбора настроек: %q", output.String())
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "выход",
			input: "2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name:  "неверный выбор",
			input: "0\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Неверный выбор\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n",
		},
		{
			name: "детерминированная победа",
			input: "1\n3\n1\n" +
				"сахар\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выберите категорию:\n 1. Животные\n 2. Города\n 3. Едa\n" +
				"Выберите сложность:\n 1. Лёгкая\n 2. Сложная\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ сахар\n" +
				"Осталось попыток: 5\n" +
				"\033[32mПобеда! Слово угадано за 1 попытки\033[0m\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name: "две игры подряд",
			input: "1\n3\n1\n" +
				"сахар\n" +
				"1\n3\n1\n" +
				"сахар\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выберите категорию:\n 1. Животные\n 2. Города\n 3. Едa\n" +
				"Выберите сложность:\n 1. Лёгкая\n 2. Сложная\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ сахар\n" +
				"Осталось попыток: 5\n" +
				"\033[32mПобеда! Слово угадано за 1 попытки\033[0m\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выберите категорию:\n 1. Животные\n 2. Города\n 3. Едa\n" +
				"Выберите сложность:\n 1. Лёгкая\n 2. Сложная\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ сахар\n" +
				"Осталось попыток: 5\n" +
				"\033[32mПобеда! Слово угадано за 1 попытки\033[0m\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name: "отображение истории попыток",
			input: "1\n3\n1\n" +
				"банан\n" +
				"сахар\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выберите категорию:\n 1. Животные\n 2. Города\n 3. Едa\n" +
				"Выберите сложность:\n 1. Лёгкая\n 2. Сложная\n" +
				"Введите слово:\n" +
				"❌✅❌✅❌ банан\n" +
				"Осталось попыток: 5\n" +
				"Введите слово:\n" +
				"❌✅❌✅❌ банан\n" +
				"✅✅✅✅✅ сахар\n" +
				"Осталось попыток: 4\n" +
				"\033[32mПобеда! Слово угадано за 2 попытки\033[0m\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"3. Статистика\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			if err := Run(input, output, 42); err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if got := output.String(); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunLoseAfterSixAttempts(t *testing.T) {
	input := strings.NewReader(
		"1\n3\n1\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"2\n",
	)

	output := &bytes.Buffer{}

	if err := Run(input, output, 42); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := output.String()

	want := "Осталось попыток: 0"
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %q", want)
	}

	want = "\033[31mНеудача. Загаданное слово: сахар\033[0m\n1. Новая игра\n"
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %q", want)
	}
}

func TestRunInvalidGuessDoesNotConsumeAttempt(t *testing.T) {
	input := strings.NewReader(
		"1\n3\n1\n" +
			"fh\n" +
			"банан\n" +
			"сахар\n" +
			"2\n",
	)

	output := &bytes.Buffer{}

	if err := Run(input, output, 42); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := output.String()

	if !strings.Contains(got, "word is not 5 letters long") {
		t.Errorf("output does not contain invalid word error")
	}

	if !strings.Contains(got, "❌✅❌✅❌ банан\nОсталось попыток: 5") {
		t.Errorf("invalid guess consumed an attempt")
	}

	if !strings.Contains(got, "Победа! Слово угадано за 2 попытки") {
		t.Errorf("unexpected number of used attempts")
	}
}

func TestPlayGameInterruptedDoesNotChangeStatistics(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"до выбора категории", ""},
		{"до выбора сложности", "3\n"},
		{"до первой попытки", "3\n1\n"},
		{"после неверного слова", "3\n1\nкот\n"},
		{"после принятой попытки", "3\n1\nбанан\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := Statistics{Wins: 2, Losses: 1}
			statistics := want
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			var output bytes.Buffer
			if err := playGame(scanner, &output, 42, &statistics); err != nil {
				t.Fatalf("playGame() error = %v", err)
			}
			if statistics != want {
				t.Errorf("statistics = %+v, want %+v", statistics, want)
			}
		})
	}
}

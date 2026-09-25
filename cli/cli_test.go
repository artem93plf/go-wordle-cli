package cli_test

import (
	"bytes"
	"fiveletters/cli"
	"strings"
	"testing"
)

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
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name:  "неверный выбор",
			input: "3\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Неверный выбор\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n",
		},
		{
			name: "детерминированная победа",
			input: "1\n" +
				"ведро\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ ведро\n" +
				"Осталось попыток: 5\n" +
				"Победа! Слово угадано за 1 попытки\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name: "две игры подряд",
			input: "1\n" +
				"ведро\n" +
				"1\n" +
				"ведро\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ ведро\n" +
				"Осталось попыток: 5\n" +
				"Победа! Слово угадано за 1 попытки\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Введите слово:\n" +
				"✅✅✅✅✅ ведро\n" +
				"Осталось попыток: 5\n" +
				"Победа! Слово угадано за 1 попытки\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
		{
			name: "отображение истории попыток",
			input: "1\n" +
				"банан\n" +
				"ведро\n" +
				"2\n",
			want: "1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Введите слово:\n" +
				"❌❌❌❌❌ банан\n" +
				"Осталось попыток: 5\n" +
				"Введите слово:\n" +
				"❌❌❌❌❌ банан\n" +
				"✅✅✅✅✅ ведро\n" +
				"Осталось попыток: 4\n" +
				"Победа! Слово угадано за 2 попытки\n" +
				"1. Новая игра\n" +
				"2. Выход\n" +
				"Выберите опцию:\n" +
				"Выход\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			output := &bytes.Buffer{}

			if err := cli.Run(input, output, 42); err != nil {
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
		"1\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"банан\n" +
			"2\n",
	)

	output := &bytes.Buffer{}

	if err := cli.Run(input, output, 42); err != nil {
		t.Fatalf("cli.Run() error = %v", err)
	}

	got := output.String()

	want := "Осталось попыток: 0"
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %q", want)
	}

	want = "Неудача. Загаданное слово: ведро"
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %q", want)
	}
}

func TestRunInvalidGuessDoesNotConsumeAttempt(t *testing.T) {
	input := strings.NewReader(
		"1\n" +
			"fh\n" +
			"банан\n" +
			"ведро\n" +
			"2\n",
	)

	output := &bytes.Buffer{}

	if err := cli.Run(input, output, 42); err != nil {
		t.Fatalf("cli.Run() error = %v", err)
	}

	got := output.String()

	if !strings.Contains(got, "word is not 5 letters long") {
		t.Errorf("output does not contain invalid word error")
	}

	if !strings.Contains(got, "❌❌❌❌❌ банан\nОсталось попыток: 5") {
		t.Errorf("invalid guess consumed an attempt")
	}

	if !strings.Contains(got, "Победа! Слово угадано за 2 попытки") {
		t.Errorf("unexpected number of used attempts")
	}
}

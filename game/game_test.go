package game

import (
	"errors"
	"slices"
	"strconv"
	"testing"
	"unicode/utf8"
)

func makeValidDictionary() []string {
	dictionary := make([]string, 0, 50)

	for i := range 50 {
		word := []rune{
			'а',
			'а',
			'а',
			rune('а' + i/32),
			rune('а' + i%32),
		}

		dictionary = append(dictionary, string(word))
	}

	return dictionary
}

func TestDefaultDictionary(t *testing.T) {
	unique := make(map[string]bool)
	for _, word := range DefaultDictionary() {
		if unique[word] {
			t.Errorf("слово %q повторяется в словаре", word)
		}
		unique[word] = true

		t.Run(word, func(t *testing.T) {
			if got := utf8.RuneCountInString(word); got != 5 {
				t.Errorf("в слове %q %d букв, ожидалось 5", word, got)
			}
			for _, letter := range word {
				if (letter < 'а' || letter > 'я') &&
					(letter < 'А' || letter > 'Я') && letter != 'ё' && letter != 'Ё' {
					t.Errorf("слово %q содержит нерусский символ %q", word, letter)
				}
			}
		})
	}
	if len(unique) < 50 {
		t.Errorf("в словаре %d уникальных слов, ожидалось минимум 50", len(unique))
	}
}

func TestIsDictionary(t *testing.T) {
	valid := makeValidDictionary()

	withDuplicate := makeValidDictionary()
	withDuplicate[49] = withDuplicate[0]

	tests := []struct {
		name       string
		dictionary []string
		want       bool
	}{
		{
			name:       "50 уникальных слов",
			dictionary: valid,
			want:       true,
		},
		{
			name:       "пустой словарь",
			dictionary: []string{},
			want:       false,
		},
		{
			name:       "меньше 50 слов в словаре",
			dictionary: valid[:49],
			want:       true,
		},
		{
			name:       "словарь имеет дубликаты",
			dictionary: withDuplicate,
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsDictionary(test.dictionary)

			if got != test.want {
				t.Errorf("IsDictionary() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsFiveLetterWord(t *testing.T) {
	tests := []struct {
		name       string
		dictionary []string
		want       bool
	}{
		{
			name:       "все слова состоят из 5 букв",
			dictionary: []string{"арбуз", "слово", "хутор"},
			want:       true,
		},
		{
			name:       "словарь имеет слово длиной меньше 5",
			dictionary: []string{"арбуз", "кот"},
			want:       false,
		},
		{
			name:       "словарь имеет слово длиной больше 5",
			dictionary: []string{"арбуз", "собака"},
			want:       false,
		},
		{
			name:       "слово в верхнем регистре состоит из 5 букв",
			dictionary: []string{"СЛОВО"},
			want:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsFiveLetterWord(test.dictionary)

			if got != test.want {
				t.Errorf("IsFiveLetterWord() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestIsRussianWord(t *testing.T) {
	tests := []struct {
		name       string
		dictionary []string
		want       bool
	}{
		{
			name:       "слово состоит из русских букв нижнего регистра",
			dictionary: []string{"арбуз"},
			want:       true,
		},
		{
			name:       "слово состоит из русских букв верхнего регистра",
			dictionary: []string{"АРБУЗ"},
			want:       true,
		},
		{
			name:       "слово содержит цифру",
			dictionary: []string{"арбу1"},
			want:       false,
		},
		{
			name:       "слово с ё",
			dictionary: []string{"ёжики"},
			want:       true,
		},
		{
			name:       "слово со спец символами",
			dictionary: []string{"сл-во"},
			want:       false,
		},
		{
			name:       "русская буква заменена на латинскую",
			dictionary: []string{"cлово"},
			want:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsRussianWord(test.dictionary)

			if got != test.want {
				t.Errorf("IsRussianWord() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestStartGameErrors(t *testing.T) {
	valid := makeValidDictionary()
	withDuplicate := append(makeValidDictionary(), valid[0])

	wrongLength := makeValidDictionary()
	wrongLength[0] = "аааа"

	withDigit := makeValidDictionary()
	withDigit[0] = "аа1аа"

	withSpecialCharacter := makeValidDictionary()
	withSpecialCharacter[0] = "аа-аа"

	tests := []struct {
		name        string
		dictionary  []string
		maxAttempts int
		wantErr     error
	}{
		{
			name:        "nil словарь",
			dictionary:  nil,
			maxAttempts: 6,
			wantErr:     ErrorInvalidDictionary,
		},
		{
			name:        "50 уникальных слов и дубликат",
			dictionary:  withDuplicate,
			maxAttempts: 6,
			wantErr:     ErrorInvalidDictionary,
		},
		{
			name:        "пустой словарь",
			dictionary:  []string{},
			maxAttempts: 6,
			wantErr:     ErrorInvalidDictionary,
		},
		{
			name:        "количество попыток равно нулю",
			dictionary:  valid,
			maxAttempts: 0,
			wantErr:     ErrorMaxAttemptsNotEqual,
		},
		{
			name:        "количество попыток отрицательное",
			dictionary:  valid,
			maxAttempts: -1,
			wantErr:     ErrorMaxAttemptsNotEqual,
		},
		{
			name:        "слово содержит не пять букв",
			dictionary:  wrongLength,
			maxAttempts: 6,
			wantErr:     ErrorWordNotFiveLetters,
		},
		{
			name:        "слово содержит цифру",
			dictionary:  withDigit,
			maxAttempts: 6,
			wantErr:     ErrorWordIsNotRussian,
		},
		{
			name:        "слово содержит специальный символ",
			dictionary:  withSpecialCharacter,
			maxAttempts: 6,
			wantErr:     ErrorWordIsNotRussian,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := StartGame(
				test.dictionary,
				test.maxAttempts,
				42,
			)

			if !errors.Is(err, test.wantErr) {
				t.Errorf(
					"StartGame() вернул ошибку %v, ожидалась %v",
					err,
					test.wantErr,
				)
			}
		})
	}
}

func TestStartGameSameSeedSameAnswer(t *testing.T) {
	dictionary := makeValidDictionary()

	tests := []struct {
		name string
		seed int64
	}{
		{
			name: "сид 42",
			seed: 42,
		},
		{
			name: "сид 100",
			seed: 100,
		},
		{
			name: "нулевой сид",
			seed: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			firstGame, err := StartGame(dictionary, 6, test.seed)
			if err != nil {
				t.Fatalf(
					"первый вызов StartGame() вернул ошибку: %v",
					err,
				)
			}

			secondGame, err := StartGame(dictionary, 6, test.seed)
			if err != nil {
				t.Fatalf(
					"второй вызов StartGame() вернул ошибку: %v",
					err,
				)
			}

			if firstGame.Answer != secondGame.Answer {
				t.Errorf(
					"при одинаковом сиде %d получены разные ответы: %q и %q",
					test.seed,
					firstGame.Answer,
					secondGame.Answer,
				)
			}
		})
	}
}

func TestStartGameInitialState(t *testing.T) {
	dictionary := makeValidDictionary()

	game, err := StartGame(dictionary, 6, 42)
	if err != nil {
		t.Fatalf(
			"StartGame() вернул неожиданную ошибку: %v",
			err,
		)
	}

	t.Run("максимальное количество попыток установлено", func(t *testing.T) {
		if game.MaxAttempts != 6 {
			t.Errorf(
				"MaxAttempts = %d, ожидалось 6",
				game.MaxAttempts,
			)
		}
	})

	t.Run("использованные попытки равны нулю", func(t *testing.T) {
		if game.UsedAttempts != 0 {
			t.Errorf(
				"UsedAttempts = %d, ожидалось 0",
				game.UsedAttempts,
			)
		}
	})

	t.Run("история игры пустая", func(t *testing.T) {
		if len(game.History) != 0 {
			t.Errorf(
				"длина History = %d, ожидалось 0",
				len(game.History),
			)
		}
	})

	t.Run("игра находится в процессе", func(t *testing.T) {
		if game.Status != StatusGameInProgress {
			t.Errorf(
				"Status = %v, ожидалось %v",
				game.Status,
				StatusGameInProgress,
			)
		}
	})

	t.Run("загаданное слово находится в словаре", func(t *testing.T) {
		if !slices.Contains(dictionary, game.Answer) {
			t.Errorf(
				"загаданное слово %q отсутствует в словаре",
				game.Answer,
			)
		}
	})
}

func TestStartGamePositiveMaxAttempts(t *testing.T) {
	tests := []struct {
		name        string
		maxAttempts int
	}{
		{name: "одна попытка", maxAttempts: 1},
		{name: "шесть попыток", maxAttempts: 6},
		{name: "десять попыток", maxAttempts: 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			game, err := StartGame(makeValidDictionary(), test.maxAttempts, 42)
			if err != nil {
				t.Fatalf("StartGame() вернул неожиданную ошибку: %v", err)
			}
			if game.MaxAttempts != test.maxAttempts {
				t.Errorf("MaxAttempts = %d, ожидалось %d", game.MaxAttempts, test.maxAttempts)
			}
		})
	}
}

func TestStartGameCopiesDictionary(t *testing.T) {
	dictionary := makeValidDictionary()

	game, err := StartGame(dictionary, 6, 42)
	if err != nil {
		t.Fatalf(
			"StartGame() вернул неожиданную ошибку: %v",
			err,
		)
	}

	if !slices.Equal(game.Dictionary, dictionary) {
		t.Fatalf("словарь сессии не совпадает с исходным по содержимому и порядку")
	}
	wantDictionary := slices.Clone(dictionary)

	// Изменяем исходный словарь после создания игры.
	for i := range dictionary {
		dictionary[i] = "яяяяя"
	}

	if !slices.Equal(game.Dictionary, wantDictionary) {
		t.Errorf(
			"словарь игры изменился после изменения исходного словаря: получено %q, ожидалось %q",
			game.Dictionary,
			wantDictionary,
		)
	}
}

func TestStartGameSmallDictionary(t *testing.T) {
	valid := makeValidDictionary()
	for _, size := range []int{1, 2, 49} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			dictionary := valid[:size]
			session, err := StartGame(dictionary, 6, 42)
			if err != nil {
				t.Fatalf("StartGame() вернул неожиданную ошибку: %v", err)
			}
			if !slices.Contains(dictionary, session.Answer) {
				t.Errorf("ответ %q отсутствует в переданном словаре", session.Answer)
			}
		})
	}
}

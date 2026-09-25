package game

import (
	"math/rand"
	"slices"
	"strings"
	"unicode"
)

// StatusGame описывает статус игровой сессии.
type StatusGame int

// LetterMark описывает результат проверки буквы.
type LetterMark string

// Возможные результаты проверки буквы.
const (
	MarkCorrect LetterMark = "✅"
	MarkPresent LetterMark = "🟡"
	MarkAbsent  LetterMark = "❌"
)

// Возможные статусы игры: игра продолжается, победа и поражение.
const (
	StatusGameInProgress StatusGame = iota
	StatusGameWin
	StatusGameLose
)

// GameSession хранит состояние одной игры и допустимый словарь.
type GameSession struct { //nolint:revive // имя задано публичным контрактом задания
	Dictionary   []string
	MaxAttempts  int
	Answer       string
	UsedAttempts int
	History      []string
	Status       StatusGame
}

// GuessResult хранит результат проверки буквы.
type GuessResult struct {
	Guess string
	Marks []LetterMark
}

// IsRussianStr проверяет, что строка содержит только русские буквы.
func IsRussianStr(word string) bool {
	for _, ch := range word {
		if !unicode.Is(unicode.Cyrillic, ch) {
			return false
		}

	}
	return true
}

// IsRussianWord проверяет, что слова содержат только русские буквы.
func IsRussianWord(dictionary []string) bool {
	for _, word := range dictionary {
		if !IsRussianStr(word) {
			return false
		}

	}
	return true
}

// IsDictionary проверяет, что словарь непустой и не содержит повторов.
func IsDictionary(dictionary []string) bool {
	unique := make(map[string]struct{})

	for _, word := range dictionary {
		if _, ok := unique[word]; ok {
			return false
		}
		unique[word] = struct{}{}
	}

	return len(unique) > 0
}

func IsFiveLetterStr(word string) bool {
	wordRune := []rune(word)
	return len(wordRune) == 5
}

// IsFiveLetterWord проверяет, что каждое слово содержит ровно пять рун.
func IsFiveLetterWord(dictionary []string) bool {
	for _, word := range dictionary {
		if !IsFiveLetterStr(word) {
			return false
		}
	}
	return true
}

// StartGame проверяет параметры и создаёт сессию с выбором ответа по seed.
func StartGame(dictionary []string, maxAttempts int, seed int64) (GameSession, error) {
	if !IsDictionary(dictionary) {
		return GameSession{}, ErrorInvalidDictionary
	}
	if maxAttempts <= 0 {
		return GameSession{}, ErrorMaxAttemptsNotEqual
	}
	if !IsFiveLetterWord(dictionary) {
		return GameSession{}, ErrorWordNotFiveLetters
	}
	if !IsRussianWord(dictionary) {
		return GameSession{}, ErrorWordIsNotRussian
	}

	dictionaryCopy := make([]string, len(dictionary))
	for i, word := range dictionary {
		dictionaryCopy[i] = strings.ToLower(word)
	}

	return GameSession{
		Dictionary:   dictionaryCopy,
		MaxAttempts:  maxAttempts,
		Answer:       dictionaryCopy[rand.New(rand.NewSource(seed)).Intn(len(dictionaryCopy))],
		UsedAttempts: 0,
		History:      []string{},
		Status:       StatusGameInProgress,
	}, nil
}

// ApplyGuess проверяет букву и отмечает ее в сессии.
func ApplyGuess(session *GameSession, guess string) (GuessResult, error) {
	if session.Status != StatusGameInProgress {
		return GuessResult{}, ErrorGameFinished
	}

	guess = strings.ToLower(guess)

	if !IsFiveLetterStr(guess) {
		return GuessResult{}, ErrorWordNotFiveLetters
	}
	if !IsRussianStr(guess) {
		return GuessResult{}, ErrorWordIsNotRussian
	}

	if !slices.Contains(session.Dictionary, guess) {
		return GuessResult{}, ErrorWordNotInDictionary
	}

	answerRune := []rune(session.Answer)
	guessRune := []rune(guess)

	marks := make([]LetterMark, len(answerRune))
	used := make([]bool, len(answerRune))

	for i := range guessRune {
		if guessRune[i] == answerRune[i] {
			marks[i] = MarkCorrect
			used[i] = true
		}
	}
	for i := range guessRune {
		if marks[i] == MarkCorrect {
			continue
		}

		marks[i] = MarkAbsent

		for j := range answerRune {
			if !used[j] && guessRune[i] == answerRune[j] {
				marks[i] = MarkPresent
				used[j] = true
				break
			}
		}
	}

	session.UsedAttempts++
	session.History = append(session.History, guess)

	if guess == session.Answer {
		session.Status = StatusGameWin
	} else if session.UsedAttempts >= session.MaxAttempts {
		session.Status = StatusGameLose
	}

	return GuessResult{
		Guess: guess,
		Marks: marks,
	}, nil
}

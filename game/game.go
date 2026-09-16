package game

import (
	"math/rand"
	"strings"
	"unicode"
)

// StatusGame описывает статус игровой сессии.
type StatusGame int

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

// IsRussianWord проверяет, что слова содержат только русские буквы.
func IsRussianWord(dictionary []string) bool {
	for _, word := range dictionary {
		for _, ch := range word {
			if !unicode.Is(unicode.Cyrillic, ch) {
				return false
			}
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

// IsFiveLetterWord проверяет, что каждое слово содержит ровно пять рун.
func IsFiveLetterWord(dictionary []string) bool {
	for _, word := range dictionary {
		wordRune := []rune(word)
		if len(wordRune) != 5 {
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

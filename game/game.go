package game

import "math/rand"

type StatusGame int

const (
	StatusGameInProgress StatusGame = iota
	StatusGameWin
	StatusGameLose
)

type GameSession struct {
	Dictionary   []string
	MaxAttempts  int
	Answer       string
	UsedAttempts int
	History      []string
	Status       StatusGame
}

func IsRussianWord(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if (ch < 'а' || ch > 'я') &&
			(ch < 'А' || ch > 'Я') &&
			ch != 'ё' && ch != 'Ё' {
			return false
		}
	}
	return true
}

func StartGame(dictionary []string, maxAttempts int, seed int64) (GameSession, error) {
	if len(dictionary) == 0 {
		return GameSession{}, ErrorDictionaryIsEmpty
	}
	if maxAttempts <= 0 {
		return GameSession{}, ErrorMaxAttemptsNotEqual
	}
	for _, word := range dictionary {
		wordRune := []rune(word)
		if len(wordRune) != 5 {
			return GameSession{}, ErrorWordNotFiveLetters
		}
		if !IsRussianWord(string(wordRune)) {
			return GameSession{}, ErrorWordIsNotRussian
		}
	}

	dictionaryCopy := make([]string, len(dictionary))
	copy(dictionaryCopy, dictionary)

	return GameSession{
		Dictionary:   dictionaryCopy,
		MaxAttempts:  maxAttempts,
		Answer:       dictionaryCopy[rand.New(rand.NewSource(seed)).Intn(len(dictionaryCopy))],
		UsedAttempts: 0,
		History:      []string{},
		Status:       StatusGameInProgress,
	}, nil
}

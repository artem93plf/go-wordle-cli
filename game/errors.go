package game

import (
	"errors"
)

var (
	// ErrorInvalidDictionary означает, что словарь пустой или содержит повторы.
	ErrorInvalidDictionary = errors.New("dictionary must be non-empty and contain no duplicates")
	// ErrorMaxAttemptsNotEqual означает, что лимит попыток неположительный.
	ErrorMaxAttemptsNotEqual = errors.New("max attempts is not equal")
	// ErrorWordNotFiveLetters означает, что длина слова отличается от пяти букв.
	ErrorWordNotFiveLetters = errors.New("word is not 5 letters long")
	// ErrorWordIsNotRussian означает, что слово содержит нерусские символы.
	ErrorWordIsNotRussian = errors.New("word is not Russian")
)

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
	// ErrorWordNotInDictionary означает, что попытки нет в словаре игры.
	ErrorWordNotInDictionary = errors.New("word is not in dictionary")
	// ErrorGameFinished означает, что попытка сделана после завершения игры.
	ErrorGameFinished = errors.New("game is already finished")
)

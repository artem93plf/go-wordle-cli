package game

import (
	"errors"
)

var (
	ErrorDictionaryIsEmpty   = errors.New("dictionary is empty")
	ErrorMaxAttemptsNotEqual = errors.New("max attempts is not equal")
	ErrorWordNotFiveLetters  = errors.New("word is not 5 letters long")
	ErrorWordIsNotRussian    = errors.New("word is not Russian")
)

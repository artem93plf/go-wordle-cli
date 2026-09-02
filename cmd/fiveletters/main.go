// Package main implements the Five Letters console game.
package main

// GameSession represents the state of one game.
type GameSession struct{}

// GuessResult represents the result of one accepted guess.
type GuessResult struct{}

// StartGame creates a new game session.
func StartGame(dictionary []string, maxAttempts int, seed int64) (GameSession, error) {
	_, _, _ = dictionary, maxAttempts, seed
	panic("not implemented")
}

// ApplyGuess validates and applies a guess to the session.
func ApplyGuess(session *GameSession, guess string) (GuessResult, error) {
	_, _ = session, guess
	panic("not implemented")
}

func main() {}

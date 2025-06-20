package domain

import (
	"log/slog"
	"time"
)

// type MessagePayload interface {
// 	TaskRequest | SolutionSubmission
// }

// type TaskRequest struct {
// 	ID       string `json:"id"`
// 	ClientID string `json:"client_id"` // TODO:
// }

// type TaskResponse struct {
// 	ID      string `json:"id"`
// 	Payload string `json:"payload"`
// }

// type SolutionSubmission struct {
// 	ID       string `json:"id"`
// 	Solution string `json:"solution"`
// }

// Quote is a struct representing a quote from Yoda
type Quote struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Source   string `json:"source"`
	Episode  string `json:"episode"`
	Year     int    `json:"year"`
	Category string `json:"category"`
}

func (quote Quote) Log() []any {
	return []any{
		slog.String("ID", quote.ID),
		slog.String("Text", quote.Text),
		slog.String("Source", quote.Source),
		slog.String("Episode", quote.Episode),
		slog.Int("Year", quote.Year),
		slog.String("Category", quote.Category),
	}
}

// Challenge contains data of the puzzle given by the server
type Challenge struct {
	ID         string    `json:"id"`         // unique identifier of the puzzle
	Data       string    `json:"data"`       // random data + timestamp
	Difficulty int       `json:"difficulty"` // number of leading zero bits in the hash
	ExpiresAt  time.Time `json:"expiresAt"`  // time of expiration
}

func (challenge Challenge) Log() []any {
	return []any{
		slog.String("ID", challenge.ID),
		slog.String("Data", challenge.Data),
		slog.Int("Difficulty", challenge.Difficulty),
		slog.Time("ExpiresAt", challenge.ExpiresAt),
	}
}

type Solution struct {
	ID    string `json:"id"`
	Nonce uint64 `json:"payload,omitempty"`
}

func (solution Solution) Log() []any {
	return []any{
		slog.String("ID", solution.ID),
		slog.Uint64("Nonce", solution.Nonce),
	}
}

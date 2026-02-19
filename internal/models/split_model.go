package models

import "time"

type Split struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	TotalAmount  float64       `json:"total_amount"`
	Currency     string        `json:"currency"`
	Strategy     SplitStrategy `json:"strategy"`
	Participants []Participant `json:"participants"`
	CreatedBy    string        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Participant struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email,omitempty"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage,omitempty"`
	IsPaid     bool    `json:"is_paid"`
}

type SplitStrategy string

const (
	StrategyEqual      SplitStrategy = "equal"
	StrategyExact      SplitStrategy = "exact"
	StrategyPercentage SplitStrategy = "percentage"
)

type CreateSplitRequest struct {
	Title        string               `json:"title"`
	TotalAmount  float64              `json:"total_amount"`
	Currency     string               `json:"currency"`
	Strategy     SplitStrategy        `json:"strategy"`
	CreatedBy    string               `json:"created_by"`
	Participants []ParticipantRequest `json:"participants"`
}

type ParticipantRequest struct {
	Name       string  `json:"name"`
	Email      string  `json:"email,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}

type SettleRequest struct {
	ParticipantID string `json:"participant_id"`
}

type SplitSummary struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	TotalAmount      float64   `json:"total_amount"`
	Currency         string    `json:"currency"`
	ParticipantCount int       `json:"participant_count"`
	SettledCount     int       `json:"settled_count"`
	CreatedAt        time.Time `json:"created_at"`
}

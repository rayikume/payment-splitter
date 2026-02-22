package services

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rayikume/payment-splitter/internal/models"
)

var (
	ErrParticipantNotFound         = errors.New("participant not found")
	ErrSplitNotFound               = errors.New("split not found")
	ErrPercentageMismatch          = errors.New("participant percentages do not sum to 100")
	ErrAmountMismatch              = errors.New("participant amounts do not sum to total")
	ErrInvalidStrategy             = errors.New("invalid split strategy")
	ErrInvalidAmount               = errors.New("total amount must be greater than 0")
	ErrLessThanMinimumParticipants = errors.New("at least 2 participants are required")
)

type SplitService struct {
	mux    sync.RWMutex
	splits map[string]*models.Split
}

func NewSplitService() *SplitService {
	return &SplitService{
		splits: make(map[string]*models.Split),
	}
}

func (s *SplitService) Create(req models.CreateSplitRequest) (*models.Split, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	participants, erro := calculateShares(req)
	if erro != nil {
		return nil, erro
	}

	now := time.Now().UTC()
	split := &models.Split{
		ID:           uuid.NewString(),
		Title:        req.Title,
		TotalAmount:  req.TotalAmount,
		Currency:     req.Currency,
		Strategy:     req.Strategy,
		Participants: participants,
		CreatedBy:    req.CreatedBy,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.mux.Lock()
	s.splits[split.ID] = split
	s.mux.Unlock()

	return split, nil
}

func (s *SplitService) Settle(splitID, participantID string) (*models.Split, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	split, ok := s.splits[splitID]
	if !ok {
		return nil, ErrSplitNotFound
	}

	found := false
	for i := range split.Participants {
		if split.Participants[i].ID == participantID {
			split.Participants[i].IsPaid = true
			found = true
			break
		}
	}

	if !found {
		return nil, ErrParticipantNotFound
	}

	split.UpdatedAt = time.Now().UTC()
	return split, nil
}

func (s *SplitService) GetByID(id string) (*models.Split, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	split, ok := s.splits[id]
	if !ok {
		return nil, ErrSplitNotFound
	}

	return split, nil
}

func (s *SplitService) Delete(id string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.splits[id]
	if !ok {
		return ErrSplitNotFound
	}

	delete(s.splits, id)
	return nil
}

func validateCreateRequest(req models.CreateSplitRequest) error {
	if req.TotalAmount <= 0 {
		return ErrInvalidAmount
	}
	if len(req.Participants) < 2 {
		return ErrLessThanMinimumParticipants
	}
	return nil
}

func calculateShares(req models.CreateSplitRequest) ([]models.Participant, error) {
	switch req.Strategy {
	case models.StrategyEqual:
		return splitEqual(req)
	case models.StrategyExact:
		return splitExact(req)
	case models.StrategyPercentage:
		return splitByPercetage(req)
	default:
		return nil, ErrInvalidStrategy
	}
}

func splitEqual(req models.CreateSplitRequest) ([]models.Participant, error) {
	n := len(req.Participants)
	share := math.Floor(req.TotalAmount/float64(n)*100) / 100
	remainder := math.Round((req.TotalAmount-share*float64(n))*100) / 100

	participants := make([]models.Participant, n)
	for i, p := range req.Participants {
		amount := share
		if i == n-1 {
			amount = math.Round((share+remainder)*100) / 100
		}
		participants[i] = models.Participant{
			ID:     uuid.NewString(),
			Name:   p.Name,
			Email:  p.Email,
			Amount: amount,
			IsPaid: false,
		}
	}
	return participants, nil
}

func splitExact(req models.CreateSplitRequest) ([]models.Participant, error) {
	var sum float64
	n := len(req.Participants)

	participants := make([]models.Participant, n)
	for i, p := range req.Participants {
		sum += p.Amount
		participants[i] = models.Participant{
			ID:     uuid.NewString(),
			Name:   p.Name,
			Email:  p.Email,
			Amount: math.Round(p.Amount*100) / 100,
			IsPaid: false,
		}
	}

	if math.Abs(sum-req.TotalAmount) > 0.01 {
		return nil, fmt.Errorf("%w: got %.2f, expected %.2f", ErrAmountMismatch, sum, req.TotalAmount)
	}
	return participants, nil
}

func splitByPercetage(req models.CreateSplitRequest) ([]models.Participant, error) {
	var totalPrcntg float64
	n := len(req.Participants)
	participants := make([]models.Participant, n)

	for i, p := range req.Participants {
		totalPrcntg += p.Percentage
		amount := math.Round(req.TotalAmount*p.Percentage) / 100
		participants[i] = models.Participant{
			ID:         uuid.NewString(),
			Name:       p.Name,
			Email:      p.Email,
			Amount:     amount,
			Percentage: p.Percentage,
			IsPaid:     false,
		}
	}

	if math.Abs(totalPrcntg-100) > 0.01 {
		return nil, fmt.Errorf("%w: got %.2f%%", ErrPercentageMismatch, totalPrcntg)
	}
	return participants, nil
}

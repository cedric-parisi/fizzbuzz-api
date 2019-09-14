package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cedric-parisi/fizzbuzz-api/internal/errors"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
)

// FizzbuzzService implements Service interface for fizzbuzz logic.
type FizzbuzzService struct {
	repo fizzbuzz.Repository
}

// NewFizzbuzzService creates a new service impl for fizzbuzz.
func NewFizzbuzzService(repo fizzbuzz.Repository) *FizzbuzzService {
	return &FizzbuzzService{
		repo: repo,
	}
}

// GetFizzbuzz computes and returns a fizzbuzz
func (s *FizzbuzzService) GetFizzbuzz(ctx context.Context, req *models.Fizzbuzz) ([]string, error) {
	// Perform business validation.
	if err := req.Validate(); err != nil {
		return nil, errors.InvalidArgument("invalid fizzbuzz request", err)
	}

	res := computeFizzbuzz(req)

	// Save fizzbuzz request to storage for stats.
	if err := s.repo.SaveFizzbuzz(ctx, req); err != nil {
		return nil, errors.Internal(err)
	}

	return res, nil
}

// GetMostAskedFizzbuzz returns the most asked fizzbuzz query.
func (s *FizzbuzzService) GetMostAskedFizzbuzz(ctx context.Context) (*models.Fizzbuzz, int, error) {
	fb, ct, err := s.repo.GetMostAskedFizzbuzz(ctx)
	if err != nil {
		if err == models.ErrNotFound {
			return nil, 0, errors.NotFound(fmt.Errorf("most asked fizzbuzz: %v", err))
		}
		return nil, 0, errors.Internal(err)
	}
	return fb, ct, nil
}

func computeFizzbuzz(req *models.Fizzbuzz) []string {
	var res []string
	for i := 1; i <= req.Limit; i++ {
		tmp := ""
		if i%req.Int1 == 0 {
			tmp += req.Str1
		}
		if i%req.Int2 == 0 {
			tmp += req.Str2
		}
		if tmp == "" {
			tmp = strconv.Itoa(i)
		}
		res = append(res, tmp)
	}
	return res
}

package fizzbuzz

import (
	"context"

	"github.com/cedric-parisi/fizzbuzz-api/models"
)

// Service provides fizzbuzz usecases methods.
type Service interface {
	GetFizzbuzz(ctx context.Context, req *models.Fizzbuzz) ([]string, error)
	GetMostAskedFizzbuzz(ctx context.Context) (*models.Fizzbuzz, int, error)
}

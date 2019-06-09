package fizzbuzz

import (
	"context"

	"github.com/cedric-parisi/fizzbuzz-api/models"
)

// Repository provides storage methods for fizzbuzz.
type Repository interface {
	SaveFizzbuzz(ctx context.Context, req *models.Fizzbuzz) error
	GetMostAskedFizzbuzz(ctx context.Context) (*models.Fizzbuzz, int, error)
}

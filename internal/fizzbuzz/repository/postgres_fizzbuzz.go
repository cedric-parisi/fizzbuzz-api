package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
)

const (
	insertFizzbuzzQuery          = `INSERT INTO stats(int1, int2, max_limit, str1, str2, occured_at) VALUES ($1, $2, $3, $4, $5, now());`
	selectMostAskedFizzbuzzQuery = `
	SELECT COUNT(*) AS hits, int1, int2, max_limit, str1, str2 
	FROM stats
	GROUP BY int1, int2, max_limit, str1, str2
	ORDER BY hits DESC
	LIMIT 1;`
)

type pqRepository struct {
	db              *sqlx.DB
	timeoutDuration time.Duration
}

// NewPostgresFizzbuzzRepository creates a new postgres impl for fizzbuzz.
func NewPostgresFizzbuzzRepository(db *sqlx.DB, t int) fizzbuzz.Repository {
	return &pqRepository{
		db:              db,
		timeoutDuration: time.Duration(t) * time.Second,
	}
}

// SaveFizzbuzz persists fizzbuzz request into storage.
func (s *pqRepository) SaveFizzbuzz(ctx context.Context, fb *models.Fizzbuzz) error {
	// Close the db call if too long
	ctx, cancel := context.WithTimeout(ctx, s.timeoutDuration)
	defer cancel()

	// Persists the query param to the DB.
	if _, err := s.db.ExecContext(ctx, insertFizzbuzzQuery, fb.Int1, fb.Int2, fb.Limit, fb.Str1, fb.Str2); err != nil {
		return err
	}

	return nil
}

// GetMostAskedFizzbuzz returns the most asked fizzbuzz query.
func (s *pqRepository) GetMostAskedFizzbuzz(ctx context.Context) (*models.Fizzbuzz, int, error) {
	// Close the db call if too long
	ctx, cancel := context.WithTimeout(ctx, s.timeoutDuration)
	defer cancel()

	// Retrieve the most asked fizzbuzz from the DB.
	var int1, int2, limit int
	var str1, str2 string
	var count int
	if err := s.db.QueryRowContext(ctx, selectMostAskedFizzbuzzQuery).Scan(&count, &int1, &int2, &limit, &str1, &str2); err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, models.ErrNotFound
		}
		return nil, 0, err
	}

	return &models.Fizzbuzz{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}, count, nil
}

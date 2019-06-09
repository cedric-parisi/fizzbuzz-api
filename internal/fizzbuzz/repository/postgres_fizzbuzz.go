package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
)

const (
	insertFizzbuzzQuery          = `INSERT INTO stats(checksum_query, occured_at) VALUES ($1, now());`
	selectMostAskedFizzbuzzQuery = `
	SELECT COUNT(*) AS max_asked, checksum_query 
	FROM stats
	GROUP BY checksum_query
	ORDER BY max_asked DESC
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
	// Encode the fizzbuzz struct to hexadecimal string.
	src := &bytes.Buffer{}
	if err := json.NewEncoder(src).Encode(fb); err != nil {
		return err
	}
	dst := make([]byte, hex.EncodedLen(len(src.Bytes())))

	// We ignore the Encode result as
	// this value is always EncodedLen(len(src))
	hex.Encode(dst, src.Bytes())

	// Close the db call if too long
	ctx, cancel := context.WithTimeout(ctx, s.timeoutDuration)
	defer cancel()

	// Persists the hex encoded query to the DB.
	if _, err := s.db.ExecContext(ctx, insertFizzbuzzQuery, string(dst)); err != nil {
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
	var checksum string
	var count int
	if err := s.db.QueryRowContext(ctx, selectMostAskedFizzbuzzQuery).Scan(&count, &checksum); err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, models.ErrNotFound
		}
		return nil, 0, err
	}

	// Decode the hexadecimal string.
	src := []byte(checksum)
	dst := make([]byte, len(src))
	hex.Decode(dst, src)

	// Decode the byte array into fizzbuzz struct.
	res := models.Fizzbuzz{}
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&res); err != nil {
		return nil, 0, err
	}

	return &res, count, nil
}

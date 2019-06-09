package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v2"

	"github.com/cedric-parisi/fizzbuzz-api/models"
	"github.com/jmoiron/sqlx"
)

func Test_pqRepository_SaveFizzbuzz(t *testing.T) {
	type args struct {
		ctx context.Context
		fb  *models.Fizzbuzz
	}
	tests := []struct {
		name      string
		args      args
		mockCalls func(m sqlmock.Sqlmock)
		wantErr   bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				fb: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 5,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO stats").WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "KO_TooLong",
			args: args{
				ctx: context.Background(),
				fb: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 5,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO stats").WillDelayFor(time.Duration(2) * time.Second)
			},
			wantErr: true,
		},
		{
			name: "KO_DBFailed",
			args: args{
				ctx: context.Background(),
				fb: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 5,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectExec("INSERT INTO stats").WillReturnError(fmt.Errorf("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			db, mock, _ := sqlmock.New()
			s := NewPostgresFizzbuzzRepository(sqlx.NewDb(db, "sqlmock"), 1)
			tt.mockCalls(mock)

			// Act & Assert
			if err := s.SaveFizzbuzz(tt.args.ctx, tt.args.fb); (err != nil) != tt.wantErr {
				t.Errorf("pqRepository.SaveFizzbuzz() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_pqRepository_GetMostAskedFizzbuzz(t *testing.T) {
	rows := sqlmock.NewRows([]string{"hits", "int1", "int2", "max_limit", "str1", "str2"})
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		mockCalls func(m sqlmock.Sqlmock)
		want      *models.Fizzbuzz
		want1     int
		wantErr   bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT COUNT").
					WillReturnRows(rows.AddRow(
						42,
						3, 5, 15, "fizz", "buzz",
					))
			},
			want: &models.Fizzbuzz{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want1: 42,
		},
		{
			name: "KO_TooLong",
			args: args{
				ctx: context.Background(),
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT COUNT").WillDelayFor(time.Duration(2) * time.Second)
			},
			wantErr: true,
		},
		{
			name: "KO_NoRows",
			args: args{
				ctx: context.Background(),
			},
			mockCalls: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT COUNT").WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			db, mock, _ := sqlmock.New()
			s := NewPostgresFizzbuzzRepository(sqlx.NewDb(db, "sqlmock"), 1)
			tt.mockCalls(mock)

			// Act
			got, got1, err := s.GetMostAskedFizzbuzz(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("pqRepository.GetMostAskedFizzbuzz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pqRepository.GetMostAskedFizzbuzz() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("pqRepository.GetMostAskedFizzbuzz() got1 = %v, want %v", got1, tt.want1)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

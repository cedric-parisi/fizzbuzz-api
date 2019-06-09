// +build !integration

package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
)

func Test_svc_GetFizzbuzz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := fizzbuzz.NewMockRepository(ctrl)

	type args struct {
		ctx context.Context
		req *models.Fizzbuzz
	}
	tests := []struct {
		name     string
		args     args
		mockCall func(m *fizzbuzz.MockRepository)
		want     []string
		wantErr  bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				req: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 5,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			mockCall: func(m *fizzbuzz.MockRepository) {
				m.
					EXPECT().
					SaveFizzbuzz(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			want: []string{
				"1",
				"2",
				"fizz",
				"4",
				"buzz",
			},
		},
		{
			name: "KO_InvalidRequest",
			args: args{
				ctx: context.Background(),
				req: &models.Fizzbuzz{
					Int1: 0,
				},
			},
			mockCall: func(m *fizzbuzz.MockRepository) {},
			wantErr:  true,
		},
		{
			name: "KO_FailedSaveFizzbuzz",
			args: args{
				ctx: context.Background(),
				req: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 5,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			mockCall: func(m *fizzbuzz.MockRepository) {
				m.
					EXPECT().
					SaveFizzbuzz(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.mockCall(m)
			s := NewFizzbuzzService(m)

			// Act
			got, err := s.GetFizzbuzz(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("svc.GetFizzbuzz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("svc.GetFizzbuzz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_svc_GetMostAskedFizzbuzz(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := fizzbuzz.NewMockRepository(ctrl)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		mockCall func(m *fizzbuzz.MockRepository)
		want     *models.Fizzbuzz
		want1    int
		wantErr  bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
			},
			mockCall: func(m *fizzbuzz.MockRepository) {
				m.
					EXPECT().
					GetMostAskedFizzbuzz(gomock.Any()).
					Return(&models.Fizzbuzz{
						Int1:  3,
						Int2:  5,
						Limit: 5,
						Str1:  "fizz",
						Str2:  "buzz",
					}, 42, nil)
			},
			want: &models.Fizzbuzz{
				Int1:  3,
				Int2:  5,
				Limit: 5,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want1: 42,
		},
		{
			name: "KO_NotFound",
			args: args{
				ctx: context.Background(),
			},
			mockCall: func(m *fizzbuzz.MockRepository) {
				m.
					EXPECT().
					GetMostAskedFizzbuzz(gomock.Any()).
					Return(nil, 0, models.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "KO_FailedGetMostAskedFizzbuzz",
			args: args{
				ctx: context.Background(),
			},
			mockCall: func(m *fizzbuzz.MockRepository) {
				m.
					EXPECT().
					GetMostAskedFizzbuzz(gomock.Any()).
					Return(nil, 0, fmt.Errorf("failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.mockCall(m)
			s := NewFizzbuzzService(m)

			// Act
			got, got1, err := s.GetMostAskedFizzbuzz(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("svc.GetMostAskedFizzbuzz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Assert
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("svc.GetMostAskedFizzbuzz() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("svc.GetMostAskedFizzbuzz() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_computeFizzbuzz(t *testing.T) {
	type args struct {
		req *models.Fizzbuzz
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "OK_WithStr1Str2",
			args: args{
				req: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 15,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			want: []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz"},
		},
		{
			name: "OK_WithStr1",
			args: args{
				req: &models.Fizzbuzz{
					Int1:  3,
					Int2:  5,
					Limit: 4,
					Str1:  "fizz",
					Str2:  "buzz",
				},
			},
			want: []string{"1", "2", "fizz", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeFizzbuzz(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeFizzbuzz() = %v, want %v", got, tt.want)
			}
		})
	}
}

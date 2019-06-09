// +build !integration

package transport

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/cedric-parisi/fizzbuzz-api/internal/errors"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
)

func Test_handler_GetFizzbuzz(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := fizzbuzz.NewMockService(ctrl)

	type args struct {
		etx echo.Context
	}
	tests := []struct {
		name     string
		args     args
		mockCall func(m *fizzbuzz.MockService)
		want     int
		wantErr  bool
	}{
		{
			name: "OK",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {
				m.
					EXPECT().
					GetFizzbuzz(gomock.Any(), gomock.Any()).
					Return([]string{}, nil)
			},
			want: http.StatusOK,
		},
		{
			name: "KO_InvalidArgumentInt1",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/?int1=should_be_a_number", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {},
			want:     http.StatusBadRequest,
		},
		{
			name: "KO_InvalidArgumentInt2",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/?int1=3&int2=should_be_a_number", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {},
			want:     http.StatusBadRequest,
		},
		{
			name: "KO_InvalidArgumentLimit",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/?int1=3&int2=5&limit=should_be_a_number", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {},
			want:     http.StatusBadRequest,
		},
		{
			name: "KO_FailedGetFizzbuzz",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {
				m.
					EXPECT().
					GetFizzbuzz(gomock.Any(), gomock.Any()).
					Return(nil, errors.Internal(fmt.Errorf("failed")))
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.mockCall(m)
			h := &handler{
				svc:       m,
				errLogger: logger,
			}

			// Act & assert
			if err := h.GetFizzbuzz(tt.args.etx); (err != nil) != tt.wantErr {
				t.Errorf("handler.GetFizzbuzz() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != tt.args.etx.Response().Status {
				t.Errorf("handler.GetFizzbuzz() expected status code = %v, got %v", tt.want, tt.args.etx.Response().Status)
			}
		})
	}
}

func Test_handler_GetStats(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := fizzbuzz.NewMockService(ctrl)

	type args struct {
		etx echo.Context
	}
	tests := []struct {
		name     string
		args     args
		mockCall func(m *fizzbuzz.MockService)
		want     int
		wantErr  bool
	}{
		{
			name: "OK",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/stats/", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {
				m.
					EXPECT().
					GetMostAskedFizzbuzz(gomock.Any()).
					Return(&models.Fizzbuzz{}, 1, nil)
			},
			want: http.StatusOK,
		},
		{
			name: "KO_NoStatsFound",
			args: args{
				etx: echo.New().NewContext(
					httptest.NewRequest(http.MethodGet, "/stats/", nil),
					httptest.NewRecorder(),
				),
			},
			mockCall: func(m *fizzbuzz.MockService) {
				m.
					EXPECT().
					GetMostAskedFizzbuzz(gomock.Any()).
					Return(nil, 0, errors.NotFound(fmt.Errorf("no stats")))
			},
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tt.mockCall(m)
			h := &handler{
				svc:       m,
				errLogger: logger,
			}

			// Act & assert
			if err := h.GetStats(tt.args.etx); (err != nil) != tt.wantErr {
				t.Errorf("handler.GetStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != tt.args.etx.Response().Status {
				t.Errorf("handler.GetStats() expected status code = %v, got %v", tt.want, tt.args.etx.Response().Status)
			}
		})
	}
}

func Test_encodeError(t *testing.T) {
	err := errors.NotFound(fmt.Errorf("some error"))
	etx := echo.New().NewContext(httptest.NewRequest("GET", "/", nil), httptest.NewRecorder())
	logger := logrus.New()
	dest := &strings.Builder{}
	logger.Out = dest

	type args struct {
		ctx       echo.Context
		err       error
		errLogger *logrus.Logger
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 ResponseError
	}{
		{
			name: "OK_FromAPIError",
			args: args{
				err: err,
			},
			want: http.StatusNotFound,
			want1: ResponseError{
				Err: err.(errors.APIError),
			},
		},
		{
			name: "OK_UncaughtError",
			args: args{
				ctx:       etx,
				err:       fmt.Errorf("not an apierror"),
				errLogger: logger,
			},
			want: http.StatusInternalServerError,
			want1: ResponseError{
				Err: errors.APIError{
					Message: fmt.Errorf("not an apierror").Error(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			got, got1 := encodeError(tt.args.ctx, tt.args.err, tt.args.errLogger)

			// Assert
			if got != tt.want {
				t.Errorf("encodeError() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("encodeError() got1 = %v, want %v", got1, tt.want1)
			}
			// Ensure log for server error
			if tt.want > 500 && "" == dest.String() {
				t.Errorf("encodeError() expect log for status code > 500")
			}
		})
	}
}

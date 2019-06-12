package transport

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/cedric-parisi/fizzbuzz-api/internal/errors"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz"
	"github.com/cedric-parisi/fizzbuzz-api/models"
	"github.com/cedric-parisi/fizzbuzz-api/pkg/logging"
)

type handler struct {
	svc       fizzbuzz.Service
	errLogger *logrus.Logger
}

// FizzbuzzResponse holds the result of the fizzbuzz request.
type FizzbuzzResponse struct {
	Result []string `json:"result"`
}

// StatsResponse holds the most requested fizzbuzz and its occurrence count.
type StatsResponse struct {
	Hits    int              `json:"hits"`
	Request *models.Fizzbuzz `json:"request"`
}

// NewFizzbuzzHandler creates an handler and attach fizzbuzz routes.
func NewFizzbuzzHandler(eg *echo.Group, svc fizzbuzz.Service, errLogger *logrus.Logger) {
	h := &handler{
		svc:       svc,
		errLogger: errLogger,
	}
	eg.GET("/", h.GetFizzbuzz)
	eg.GET("/stats/", h.GetStats)
}

// GetFizzbuzz get a fizzbuzz.
func (h *handler) GetFizzbuzz(etx echo.Context) error {
	ctx := context.Background()
	if etx.Request().Context() != nil {
		ctx = etx.Request().Context()
	}

	// Perform type validation from the request.
	var err error
	var int1, int2, limit int
	if int1, err = strconv.Atoi(etx.QueryParam("int1")); err != nil {
		return etx.JSON(encodeError(etx, errors.InvalidArgument("invalid fizzbuzz request", []errors.FieldError{
			{
				Name:    "int1",
				Message: err.Error(),
			},
		}), h.errLogger))
	}
	if int2, err = strconv.Atoi(etx.QueryParam("int2")); err != nil {
		return etx.JSON(encodeError(etx, errors.InvalidArgument("invalid fizzbuzz request", []errors.FieldError{
			{
				Name:    "int2",
				Message: err.Error(),
			},
		}), h.errLogger))
	}
	if limit, err = strconv.Atoi(etx.QueryParam("limit")); err != nil {
		return etx.JSON(encodeError(etx, errors.InvalidArgument("invalid fizzbuzz request", []errors.FieldError{
			{
				Name:    "limit",
				Message: err.Error(),
			},
		}), h.errLogger))
	}

	req := &models.Fizzbuzz{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  etx.QueryParam("str1"),
		Str2:  etx.QueryParam("str2"),
	}
	fb, err := h.svc.GetFizzbuzz(ctx, req)
	if err != nil {
		return etx.JSON(encodeError(etx, err, h.errLogger))
	}

	return etx.JSON(http.StatusOK, FizzbuzzResponse{
		Result: fb,
	})
}

// GetStats get stats.
func (h *handler) GetStats(etx echo.Context) error {
	ctx := context.Background()
	if etx.Request().Context() != nil {
		ctx = etx.Request().Context()
	}

	fb, ct, err := h.svc.GetMostAskedFizzbuzz(ctx)
	if err != nil {
		return etx.JSON(encodeError(etx, err, h.errLogger))
	}
	return etx.JSON(http.StatusOK, StatsResponse{
		Request: fb,
		Hits:    ct,
	})
}

func encodeError(ctx echo.Context, err error, errLogger *logrus.Logger) (int, ResponseError) {
	sc := http.StatusInternalServerError
	var res ResponseError
	if werr, ok := err.(errors.APIError); ok {
		sc = werr.StatusCode()
		res = ResponseError{
			Err: werr,
		}
	} else {
		res = ResponseError{
			Err: errors.APIError{Message: err.Error()},
		}
	}

	// We don't care about client side errors.
	if sc >= 500 {
		logging.Log(ctx, errLogger, err)
	}
	return sc, res
}

// ResponseError wraps the error.
type ResponseError struct {
	Err errors.APIError `json:"error"`
}

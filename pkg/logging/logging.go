package logging

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Log logs an error with metadata from context.
func Log(ctx echo.Context, errLogger *logrus.Logger, err error) {
	traceID := ""
	if t := ctx.Get("trace_id"); t != nil {
		traceID = fmt.Sprintf("%v", t)
	}
	errLogger.WithFields(logrus.Fields{
		"trace_id":        traceID,
		"http.uri":        ctx.Request().RequestURI,
		"http.path":       ctx.Path(),
		"http.method":     ctx.Request().Method,
		"http.user_agent": ctx.Request().UserAgent(),
		"http.status":     ctx.Response().Status,
	}).Error(err)
}

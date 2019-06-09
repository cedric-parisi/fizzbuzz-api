package tracing

import (
	"github.com/labstack/echo/v4"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
)

// Tracing creates a new span or retrieve from context
// and add request information.
func Tracing(operationName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := opentracing.SpanFromContext(c.Request().Context())
			if span == nil {
				span = opentracing.StartSpan(operationName)
			}
			defer span.Finish()

			traceID := ""
			switch sp := span.Context().(type) {
			case jaeger.SpanContext:
				traceID = sp.TraceID().String()
			}
			c.Set("trace_id", traceID)

			span.SetTag("component", operationName)
			span.SetTag("span.kind", "server")
			span.SetTag("http.url", c.Request().Host+c.Request().RequestURI)
			span.SetTag("http.method", c.Request().Method)

			if err := next(c); err != nil {
				span.SetTag("error", true)
				c.Error(err)
			}

			span.SetTag("http.status_code", c.Response().Status)
			if c.Response().Status < 200 || c.Response().Status > 299 {
				span.SetTag("error", true)
			} else {
				span.SetTag("error", false)
			}
			return nil
		}
	}
}

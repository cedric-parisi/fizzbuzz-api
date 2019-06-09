package instrumenting

import (
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpSuccessRegex, _ = regexp.Compile("^2[0-9]{2}$")

	// HTTPRequestsTotalCounter represents a prometheus counter for counting http calls
	HTTPRequestsTotalCounter *prometheus.CounterVec

	// HTTPRequestDurationHistogram represents a promtheus histogram for measuring http calls durations
	HTTPRequestDurationHistogram *prometheus.HistogramVec
)

func init() {
	HTTPRequestsTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests received.",
	}, []string{"component", "path", "code", "method", "success"})

	HTTPRequestDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration in seconds",
	}, []string{"component", "path", "success"})

	prometheus.MustRegister(HTTPRequestsTotalCounter, HTTPRequestDurationHistogram)
}

// Metrics returns a middleware which increment the request counter
// and add the request response time to the metrics.
func Metrics(component string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func(begin time.Time) {
				status := c.Response().Status
				success := httpSuccessRegex.MatchString(strconv.Itoa(status))
				HTTPRequestsTotalCounter.With(prometheus.Labels{
					"component": component,
					"path":      c.Path(),
					"code":      strconv.Itoa(status),
					"method":    c.Request().Method,
					"success":   strconv.FormatBool(success),
				})

				HTTPRequestDurationHistogram.With(prometheus.Labels{
					"component": component,
					"path":      c.Path(),
					"success":   strconv.FormatBool(success),
				}).Observe(time.Since(begin).Seconds())

			}(time.Now())
			return next(c)
		}
	}
}

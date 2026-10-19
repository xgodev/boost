package elasticsearch

import (
	"io"
	"net/http"
	"time"

	"github.com/xgodev/boost/wrapper/log"
)

// Logger implements the elastictransport.Logger interface on the boost logger.
type Logger struct {
}

// LogRoundTrip logs one request at the level of its outcome: error for a
// transport error or a 5xx, warn for a 4xx, debug for a success.
func (l *Logger) LogRoundTrip(
	req *http.Request,
	res *http.Response,
	err error,
	start time.Time,
	dur time.Duration,
) error {
	fields := log.Fields{
		"method":    req.Method,
		"duration":  dur,
		"req_bytes": drain(req.Body),
		"url":       req.URL.String(),
	}
	if res != nil {
		fields["status_code"] = res.StatusCode
		fields["res_bytes"] = drain(res.Body)
	}

	logger := log.FromContext(req.Context()).WithFields(fields)

	switch {
	case err != nil:
		logger.Error(err)
	case res.StatusCode < 300:
		logger.Debug("success")
	case res.StatusCode < 500:
		logger.Warn("warning")
	default:
		logger.Error("error")
	}

	return nil
}

// drain counts the bytes of a body copy the transport handed to the logger.
func drain(body io.ReadCloser) int64 {
	if body == nil || body == http.NoBody {
		return 0
	}
	n, _ := io.Copy(io.Discard, body) // intentional: a short copy only undercounts a log field
	return n
}

// RequestBodyEnabled makes the client pass request body to logger.
func (l *Logger) RequestBodyEnabled() bool { return true }

// ResponseBodyEnabled makes the client pass response body to logger.
func (l *Logger) ResponseBodyEnabled() bool { return true }

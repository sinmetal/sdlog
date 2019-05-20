package buildlog

import (
	"encoding/json"
	"time"

	"go.uber.org/zap/zapcore"
)

var _ json.Marshaler = Severity(0)

// Severity provides log levels.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity
type Severity int

const (
	// SeverityDefault provides severity level.
	SeverityDefault Severity = 0
	// SeverityDebug provides severity level.
	SeverityDebug = 100
	// SeverityInfo provides severity level.
	SeverityInfo = 200
	// SeverityNotice provides severity level.
	SeverityNotice = 300
	// SeverityWarning provides severity level.
	SeverityWarning = 400
	// SeverityError provides severity level.
	SeverityError = 500
	// SeverityCritical provides severity level.
	SeverityCritical = 600
	// SeverityAlert provides severity level.
	SeverityAlert = 700
	// SeverityEmergency provides severity level.
	SeverityEmergency = 800
)

// MarshalJSON convert raw value to JSON value.
func (severity Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(severity.String())
}

// String returns Severity about string format.
func (severity Severity) String() string {
	switch severity {
	case 0:
		return "DEFAULT"
	case 100:
		return "DEBUG"
	case 200:
		return "INFO"
	case 300:
		return "NOTICE"
	case 400:
		return "WARNING"
	case 500:
		return "ERROR"
	case 600:
		return "CRITICAL"
	case 700:
		return "ALERT"
	case 800:
		return "EMERGENCY"
	default:
		return "ERROR"
	}
}

// LogEntry provides special fields in structured log.
// spec: https://cloud.google.com/logging/docs/agent/configuration#special-fields
type LogEntry struct {
	Severity       Severity                `json:"severity"`
	HTTPRequest    *HTTPRequest            `json:"httpRequest,omitempty"`
	Time           Time                    `json:"time,omitempty"`
	Trace          string                  `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string                  `json:"logging.googleapis.com/spanId,omitempty"`
	Operation      *LogEntryOperation      `json:"logging.googleapis.com/operation,omitempty"`
	SourceLocation *LogEntrySourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Message        string                  `json:"message,omitempty"`
}

// HTTPRequest provides HTTPRequest log.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
type HTTPRequest struct {
	RequestMethod                  string   `json:"requestMethod"`
	RequestURL                     string   `json:"requestUrl"`
	RequestSize                    int64    `json:"requestSize,string,omitempty"`
	Status                         int      `json:"status"`
	ResponseSize                   int64    `json:"responseSize,string,omitempty"`
	UserAgent                      string   `json:"userAgent,omitempty"`
	RemoteIP                       string   `json:"remoteIp,omitempty"`
	Referer                        string   `json:"referer,omitempty"`
	Latency                        Duration `json:"latency,omitempty"`
	CacheLookup                    *bool    `json:"cacheLookup,omitempty"`
	CacheHit                       *bool    `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer *bool    `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 *int64   `json:"cacheFillBytes,string,omitempty"`
	Protocol                       string   `json:"protocol"`
}

// LogEntryOperation provides information for long-running operation.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logentryoperation
type LogEntryOperation struct {
	ID       string `json:"id,omitempty"`
	Producer string `json:"producer,omitempty"`
	First    *bool  `json:"first,omitempty"`
	Last     *bool  `json:"last,omitempty"`
}

// LogEntrySourceLocation provides source location of log emitting.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logentrysourcelocation
type LogEntrySourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     int64  `json:"line,string,omitempty"`
	Function string `json:"function,omitempty"`
}

// MarshalLogObject with encoder.
func (sl *LogEntrySourceLocation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("file", sl.File)
	encoder.AddInt64("line", sl.Line)
	encoder.AddString("function", sl.Function)

	return nil
}

var _ json.Marshaler = Time(time.Time{})

// Time provides time.Time by protobuf format.
type Time time.Time

// MarshalJSON convert raw value to JSON value.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339Nano))
}

var _ json.Marshaler = Duration(0)

// Duration provides time.Duration by protobuf format.
type Duration time.Duration

// MarshalJSON convert raw value to JSON value.
func (d Duration) MarshalJSON() ([]byte, error) {
	nanos := time.Duration(d).Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	v := make(map[string]interface{})
	v["seconds"] = int64(secs)
	v["nanos"] = int32(nanos)

	return json.Marshal(v)
}

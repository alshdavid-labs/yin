package yin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

var color = struct {
	green   string
	white   string
	yellow  string
	red     string
	blue    string
	magenta string
	cyan    string
	reset   string
}{
	green:   string([]byte{27, 91, 57, 55, 59, 52, 50, 109}),
	white:   string([]byte{27, 91, 57, 48, 59, 52, 55, 109}),
	yellow:  string([]byte{27, 91, 57, 48, 59, 52, 51, 109}),
	red:     string([]byte{27, 91, 57, 55, 59, 52, 49, 109}),
	blue:    string([]byte{27, 91, 57, 55, 59, 52, 52, 109}),
	magenta: string([]byte{27, 91, 57, 55, 59, 52, 53, 109}),
	cyan:    string([]byte{27, 91, 57, 55, 59, 52, 54, 109}),
	reset:   string([]byte{27, 91, 48, 109}),
}

func statusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return color.green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return color.white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return color.yellow
	default:
		return color.red
	}
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return color.blue
	case "POST":
		return color.cyan
	case "PUT":
		return color.yellow
	case "DELETE":
		return color.red
	case "PATCH":
		return color.green
	case "HEAD":
		return color.magenta
	case "OPTIONS":
		return color.white
	default:
		return color.reset
	}
}

type LoggerValues struct {
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
}

type LoggerConfig struct {
	SkipPaths      []string
	NoColor        bool
	HideTimeStamp  bool
	HideStatusCode bool
	HideLatency    bool
	HideClientIP   bool
	HideMethod     bool
	HidePath       bool
}

func CreateLog(w io.Writer, values *LoggerValues, config *LoggerConfig) {
	var output string
	if config.HideTimeStamp == false {
		timeStamp := values.TimeStamp.Format("02/01/2006 - 15:04:05")
		output += fmt.Sprintf("%v | ", timeStamp)
	}
	if config.HideLatency == false {
		output += fmt.Sprintf("%13v | ", values.Latency)
	}
	if config.HideStatusCode == false {
		c := ""
		if config.NoColor == false {
			c = statusCodeColor(values.StatusCode)
		}
		output += fmt.Sprintf("%s %3d %s", c, values.StatusCode, color.reset)
	}
	if config.HideMethod == false {
		c := ""
		if config.NoColor == false {
			c = methodColor(values.Method)
		}
		output += fmt.Sprintf("%s %-7s %s", c, values.Method, color.reset)
	}
	if config.HideClientIP == false {
		output += fmt.Sprintf(" %15s", values.ClientIP)
	}
	if config.HidePath == false {
		output += fmt.Sprintf(" %s\n", values.Path)
	}
	output += fmt.Sprintf("%s", values.ErrorMessage)

	fmt.Fprint(w, output)
}

func hasColorSupport(w io.Writer) bool {
	out := w.(*os.File)
	return !(os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(out.Fd()) && !isatty.IsCygwinTerminal(out.Fd())))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logger(out io.Writer, config *LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		out := os.Stdout
		supportsColor := hasColorSupport(out)
		if supportsColor == false {
			config.NoColor = true
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			forward := ClientIP(r)
			startTime := time.Now()
			sr := &statusRecorder{w, 200}
			next.ServeHTTP(sr, r)
			latency := time.Now().Sub(startTime)
			for _, path := range config.SkipPaths {
				if strings.HasPrefix(r.URL.String(), path) {
					return
				}
			}
			log := &LoggerValues{
				TimeStamp:    time.Now(),
				StatusCode:   sr.status,
				Latency:      latency,
				ClientIP:     forward,
				Method:       r.Method,
				Path:         r.URL.String(),
				ErrorMessage: "",
			}
			CreateLog(out, log, config)
		})
	}
}

func SimpleLogger(next http.Handler) http.Handler {
	out := os.Stdout
	return Logger(out, &LoggerConfig{
		SkipPaths:     []string{"/ping"},
		HideLatency:   true,
		HideTimeStamp: true,
		HideClientIP:  true,
	})(next)
}

func DefaultLogger(next http.Handler) http.Handler {
	out := os.Stdout
	return Logger(out, &LoggerConfig{
		SkipPaths: []string{"/ping"},
	})(next)
}

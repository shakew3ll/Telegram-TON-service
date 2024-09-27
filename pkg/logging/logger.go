package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/shakew3ll/Telegram-TON-service.git/config"
)

// writerHook implements the interface logrus.Hook
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire processes the log record and writes it to the specified Writer (logrus.Hook's method)
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	text := []byte(line)

	for _, w := range hook.Writer {
		if _, err := w.Write(text); err != nil {
			return err
		}
	}

	return nil
}

// Levels returns the logging levels for which this hook is active
func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Logger wrap over logrus.Entry
type Logger struct {
	*logrus.Entry
}

// NewLogger creates a new Logger with the specified parameters
func New(cfg config.Config) (*Logger, error) {
	l := logrus.New()

	writer := cfg.Logger.Writer
	if writer == "" {
		writer = "stdout"
	}
	logLevel := cfg.Logger.Level
	if logLevel == "" {
		logLevel = "trace"
	}

	l.SetReportCaller(true)
	l.Formatter = &ColoredTextFormatter{
		TextFormatter: logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
			DisableColors: false,
			FullTimestamp: true,
		},
	}

	var output io.Writer
	switch writer {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		output = os.Stdout
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{output},
		LogLevels: logrus.AllLevels,
	})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	l.SetLevel(level)

	return &Logger{logrus.NewEntry(l)}, nil
}

// GetLoggerWithField returns the Logger with the added field
func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

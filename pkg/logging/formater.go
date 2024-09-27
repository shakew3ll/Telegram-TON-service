package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type ColoredTextFormatter struct {
	logrus.TextFormatter
}

var standardColor = "\033[0m"

const (
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
)

func (f *ColoredTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	line, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	levelStart := strings.Index(string(line), entry.Level.String())
	levelEnd := levelStart + len(entry.Level.String())

	switch entry.Level {
	case logrus.InfoLevel:
		standardColor = green
	case logrus.WarnLevel:
		standardColor = yellow
	case logrus.TraceLevel:
		standardColor = blue
	case logrus.ErrorLevel:
		standardColor = red
	case logrus.FatalLevel:
		standardColor = red
	}

	coloredLine := string(line[:levelStart]) + standardColor + string(line[levelStart:levelEnd]) + "\033[0m" + string(line[levelEnd:])

	return []byte(coloredLine), nil
}

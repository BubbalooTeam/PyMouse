package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	DisableColors bool
}

type fileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
}

const (
	TimeLogFormatter = "02/01/2006 15:04:05"
)

func (hook *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor string
	if !f.DisableColors {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = "\033[36m"
		case logrus.InfoLevel:
			levelColor = "\033[32m"
		case logrus.WarnLevel:
			levelColor = "\033[33m"
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = "\033[31m"
		default:
			levelColor = "\033[37m"
		}
	}

	log := fmt.Sprintf(
		"%s%s\033[0m | %s | %s:%s | %s | %s\n",
		levelColor,
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format(TimeLogFormatter),
		fmt.Sprintf("\033[35m%s\033[0m", entry.Caller.File),
		fmt.Sprintf("\033[35m%d\033[0m", entry.Caller.Line),
		fmt.Sprintf("\033[34m%s\033[0m", entry.Caller.Function),
		entry.Message,
	)
	if f.DisableColors {
		log = fmt.Sprintf(
			"%s | %s | %s:%d | %s | %s\n",
			strings.ToUpper(entry.Level.String()),
			entry.Time.Format(TimeLogFormatter),
			entry.Caller.File,
			entry.Caller.Line,
			entry.Caller.Function,
			entry.Message,
		)
	}
	return []byte(log), nil
}

func ConfigureLogging() error {
	file, err := os.OpenFile("PyMouse.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Errorf("Error in opening the Log file: %v", err)
		return err
	}

	// Add Logrus Hooks
	logrus.AddHook(&fileHook{Writer: file, Formatter: &CustomFormatter{DisableColors: true}})
	logrus.AddHook(&fileHook{Writer: os.Stdout, Formatter: &CustomFormatter{DisableColors: false}})

	// Setting Logrus configuration
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&CustomFormatter{})
	return nil
}

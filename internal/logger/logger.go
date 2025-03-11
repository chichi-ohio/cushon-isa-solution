package logger

import (
	"cushion-isa/internal/config"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init(cfg *config.LoggerConfig) {
	// Set logger output
	log.SetOutput(os.Stdout)

	// Set logger format
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func InitLogger(level string) {
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set global logger
	logrus.SetOutput(log.Out)
	logrus.SetFormatter(log.Formatter)
	logrus.SetLevel(log.Level)
}

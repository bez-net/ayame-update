package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
)

// for the logrus.logger, set log rotation, directory, file name
func setupLogger() *logrus.Logger {
	level, err := logrus.ParseLevel(Options.LogLevel)
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	logPath := fmt.Sprintf("%s/%s", Options.LogDir, Options.LogName)
	path, err := filepath.Abs(logPath + ".%Y%m%d")
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	rl, err := rotatelogs.New(path,
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(3600*time.Second),
	)
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	out := io.MultiWriter(os.Stdout, rl)
	logger := logrus.Logger{
		Formatter: &logrus.TextFormatter{},
		Level:     level,
		Out:       out,
	}
	logger.Info("Setup log finished.")

	return &logger
}

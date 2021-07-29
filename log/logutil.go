package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func InitLog(logFileName string) bool {

	Environment := "development"
	// do something here to set environment depending on an environment variable
	// or command-line flag
	if Environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		formatter := &logrus.TextFormatter{}
		formatter.ForceColors = true
		formatter.FullTimestamp = true

		logger.SetFormatter(formatter)
	}

	// logger.SetReportCaller(true)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//   logger.SetOutput(os.Stdout)
	if logFileName == "" {
		logFileName = "./log/gossh.log"
	}

	flag := os.O_RDWR | os.O_APPEND | os.O_CREATE
	chmod := os.FileMode(0755)

	logFile, err := os.OpenFile(logFileName, flag, chmod)
	if err != nil {
		logger.Error("unable to write file on filehook %v", err)
		return false
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	// Logging Method Name
	// If you wish to add the calling method as a field, instruct the log via:
	// logger.SetReportCaller(true)
	logger.SetOutput(mw)
	// logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logrus.InfoLevel)
	//   logger.SetLevel(logger.WarnLevel)
	return true
}

// 封装logrus.Fields
type Fields logrus.Fields

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// +++++++++++++++++ligq start

func Trace(format string, args ...interface{}) {
	logger.Logf(logrus.TraceLevel, format, args...)
}

func Debug(format string, args ...interface{}) {
	logger.Logf(logrus.DebugLevel, format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Logf(logrus.InfoLevel, format, args...)
}

func Print(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Logf(logrus.WarnLevel, format, args...)
}

func Warning(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Logf(logrus.ErrorLevel, format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Logf(logrus.FatalLevel, format, args...)
	logger.Exit(1)
}

func Panic(format string, args ...interface{}) {
	logger.Logf(logrus.PanicLevel, format, args...)
}

func SetLevel(plogLevel string) {
	switch plogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

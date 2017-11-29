package setup

import (
	logger "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

// Sets up the logrus logger
// If log path in config is nil, logs to std
// If loglevel is neither of "Debug", "Info", "Warn", "Error", e.g. "Nil" logging is disabled
func Log(env *Env) {
	if env.Log.Path != "" {
		logPath := path.Join(env.Root, env.Log.Path)
		logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatalf("error opening log file: %v", err)
		}
		logger.SetOutput(logFile)
	} else {
		logger.SetOutput(os.Stdout)
	}
	switch env.Log.Level {
	case "Debug":
		logger.SetLevel(logger.DebugLevel)
	case "Info":
		logger.SetLevel(logger.InfoLevel)
	case "Warn":
		logger.SetLevel(logger.WarnLevel)
	case "Error":
		logger.SetLevel(logger.ErrorLevel)
	default:
		logger.SetLevel(logger.ErrorLevel)
		logger.SetOutput(ioutil.Discard)
	}
	logger.SetFormatter(&logger.TextFormatter{FullTimestamp: true})
}

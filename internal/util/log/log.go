package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"equiprent/internal/util/config"
	"equiprent/internal/util/flags"

	"github.com/sirupsen/logrus"
)

// Logger is the logger for the proxymanager application
var Logger *logrus.Logger

// Initialize Logger for the cacheGrab package
func Initialize() {
	Logger = logrus.New()
	// set report caller to true
	Logger.ReportCaller = true
	formatter := &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", path.Join(path.Base(path.Dir(f.File)), path.Base(f.File)), f.Line)
		},
	}
	// dev mode logs at debug level to stdout without removing whitespace
	if *flags.Dev {
		Logger.Level = logrus.DebugLevel
		Logger.Out = os.Stdout
		formatter.PrettyPrint = true
	} else {
		// otherwise we output to a log file specified by the user
		// TODO: sns would be nice one day
		if config.Conf.DB.LogLocation == "" {
			config.Conf.DB.LogLocation = "./cacheGrab.log"
		} else {
			info, err := os.Stat(config.Conf.DB.LogLocation)
			if err != nil {
				panic(err.Error())
			}
			if info.IsDir() {
				config.Conf.DB.LogLocation = filepath.Join(config.Conf.DB.LogLocation, "cacheGrab.log")
			}
		}
		f, err := os.OpenFile(config.Conf.DB.LogLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
		Logger.Out = f
	}
	// we do this down here to control the pretty print
	Logger.SetFormatter(formatter)
	return
}

// Stop Logger
func Stop() {
	// closes the file handler on graceful shutdown
	Logger.Exit(0)
}

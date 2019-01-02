package service

import (
	"os"
	"github.com/kataras/golog"
	"reflect"
	"fmt"
	"reverse_proxy/util"
)

var (
	webLogger     *golog.Logger
	output        map[string]*os.File
)

func init() {
	output = make(map[string]*os.File)
}

func GetWebLoggerInstance() *golog.Logger {
	if (webLogger == nil) {
		webLogger = golog.New()
		logConfig := LoadWebLogConfig()
		outputFile := GetLogFile(logConfig.FilePath)
		webLogger.SetOutput(outputFile)
		webLogger.SetLevel(logConfig.Level)
		webLogger.SetTimeFormat("2006-01-02 15:04:05")
	}

	return webLogger
}

func GetLogFile(fileName string) *os.File {
	if _, ok := output[fileName]; !ok {
		outputFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		output[fileName] = outputFile
	}

	return output[fileName]
}

func LazyInfoLog(name string, lazy func() string) {
	if logInstance := GetWebLoggerInstance(); logInstance.Level >= golog.InfoLevel {
		logInstance.Info(fmt.Sprintf("routine-id:%v %s:%v", util.GetGID(), name, lazy()))
	}
}

func InfoLog(name, msg interface{}) {
	if logInstance := GetWebLoggerInstance(); logInstance.Level >= golog.InfoLevel {
		msgType := reflect.TypeOf(msg).Kind()
		switch {
		case msgType == reflect.Map:
			logInstance.Info(fmt.Sprintf("routine-id:%v %s:%s", util.GetGID(), name, util.MapToString(msg.(map[string]interface{}))))
		default:
			logInstance.Info(fmt.Sprintf("routine-id:%v %s:%v", util.GetGID(), name, msg))
		}
	}
}

func ErrorLog(name, msg interface{}) {
	if logInstance := GetWebLoggerInstance(); logInstance.Level >= golog.InfoLevel {
		msgType := reflect.TypeOf(msg).Kind()
		switch {
		case msgType == reflect.Map:
			logInstance.Error(fmt.Sprintf("routine-id:%v %s:%s\n", util.GetGID(), name, util.MapToString(msg.(map[string]interface{}))))
		default:
			logInstance.Error(fmt.Sprintf("routine-id:%v %s:%v\n", util.GetGID(), name, msg))
		}
	}
}

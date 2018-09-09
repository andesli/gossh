package logs

import (
	"github.com/astaxie/beego/logs"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

// Legacy loglevel constants to ensure backwards compatibility.
//
// Deprecated: will be removed in 1.5.0.
const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

//declare only one log instance. all the app use it to write logs
var (
	log = logs.NewLogger(1000)
)

//return  a point to log instance for other package
// all other packages in this app use this functions to get a log instance for writing logs
func NewLogger() *logs.BeeLogger {
	return log
}

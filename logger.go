package systemservice

import "log"

/*
Logger implements a basic, overridable logging interface
*/
type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type lgr struct{}

/*
Log implements the log.Println interface
*/
func (lgr) Log(v ...interface{}) {
	log.Println(v...)
}

/*
Logf implements the log.Printf interface
*/
func (lgr) Logf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

var logger Logger = lgr{}

/*
SetLogger allows the consumer of this package (that's you!) configure your
own customer logger. As long as it implements the "Logger" interface
*/
func SetLogger(customLogger Logger) {
	logger = customLogger
}

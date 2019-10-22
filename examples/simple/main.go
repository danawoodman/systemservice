package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/danawoodman/systemservice"
)

var serv systemservice.SystemService
var logger systemservice.Logger = customLogger{}

func main() {
	systemservice.SetLogger(logger)

	os := runtime.GOOS

	// Construct the command to call. On Windows,
	// we have to use a non-unix equivilent and
	// pass the absolute path to the .exe.
	var prog = "sleep"
	if os == "windows" {
		prog = "C:\\Windows\\System32\\timeout.exe"
	}

	cmd := systemservice.ServiceCommand{
		Name:          "MyService",
		Label:         "com.myservice",
		Program:       prog,
		Args:          []string{"60"},
		Description:   "My systemservice test!",
		Documentation: "https://github.com/danawoodman/systemservice",
	}

	logger.Log("created command: ", cmd.String())

	serv = systemservice.New(cmd)

	logStatus()

	if serv.Exists() {
		logger.Log("service exists, uninstalling...")

		if err := serv.Uninstall(); err != nil {
			panic(err)
		}

		logger.Log("service uninstalled")
	}

	logger.Log("running command: ", serv.Command.String())

	logStatus()

	logger.Log("installing service")

	if err := serv.Install(true); err != nil {
		panic(err)
	}

	logStatus()

	logger.Log("stopping service")

	if err := serv.Stop(); err != nil {
		panic(err)
	}

	logStatus()

	logger.Log("starting service")

	if err := serv.Start(); err != nil {
		panic(err)
	}

	logStatus()

	logger.Log("restarting service")

	if err := serv.Restart(); err != nil {
		panic(err)
	}

	logStatus()

	logger.Log("uninstall service")

	if err := serv.Uninstall(); err != nil {
		panic(err)
	}

	logStatus()
}

func logStatus() {
	stat, err := serv.Status()

	if err != nil {
		panic(err)
	}

	logger.Logf("[STATUS] running: %t, pid: %d\n", stat.Running, stat.PID)
}

// Setup a custom logger to use
type customLogger struct{}

func (customLogger) Log(v ...interface{}) {
	log.Println("[EXAMPLE] ", fmt.Sprint(v...))
}

func (customLogger) Logf(format string, v ...interface{}) {
	log.Printf("[EXAMPLE] "+format, v...)
}

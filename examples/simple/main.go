package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/danawoodman/systemservice"
)

var logger systemservice.Logger = customLogger{}
var service systemservice.SystemService

func init() {

	// Setup custom logger
	systemservice.SetLogger(logger)

	exe := exePath()
	logger.Log("executable path: ", exe)

	// Configure the process to run, in this case it is this
	// package itself but it could be any pacakge that implements
	// a valid system service
	cmd := systemservice.ServiceCommand{
		Name:          "MyService",
		Label:         "com.myservice",
		Program:       exe,
		Args:          []string{"run"},
		Description:   "My systemservice test!",
		Documentation: "https://github.com/danawoodman/systemservice",
		Debug:         true,
	}

	logger.Log("created command: ", cmd.String())

	// Configure the service with the given
	service = systemservice.New(cmd)
}

func main() {

	if len(os.Args) < 2 {
		logger.Log("no command specified")
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "run":
		run()
	case "install":
		install(true)
	case "start":
		start()
	case "stop":
		stop()
	case "uninstall":
		uninstall()
	case "status":
		status()
	}
}

func run() {
	logger.Log("Running service...")

	err := service.Run()
	if err != nil {
		panic(err)
	}
}

func status() {
	stat, err := service.Status()

	if err != nil {
		panic(err)
	}

	logger.Logf("[STATUS] running: %t, pid: %d\n", stat.Running, stat.PID)
}

func install(start bool) {
	logger.Log("[INSTALL] installing service")

	err := service.Install(start)

	if err != nil {
		panic(err)
	}

	logger.Log("[INSTALL] service installed!")
}

func uninstall() {
	logger.Log("[UNINSTALL] uninstalling service")

	err := service.Uninstall()

	if err != nil {
		panic(err)
	}

	logger.Log("[UNINSTALL] service uninstalled")
}

func start() {
	logger.Log("[START] starting service")

	err := service.Start()

	if err != nil {
		panic(err)
	}

	logger.Log("[START] service started!")
}

func stop() {
	logger.Log("[STOP] stopping service")

	err := service.Stop()

	if err != nil {
		panic(err)
	}

	logger.Log("[STOP] service stopped!")
}

/*
Setup a custom logger to override the default logger that systemservice provides (optional)
*/
type customLogger struct{}

func (customLogger) Log(v ...interface{}) {
	log.Println("[EXAMPLE] ", fmt.Sprint(v...))
}

func (customLogger) Logf(format string, v ...interface{}) {
	log.Printf("[EXAMPLE] "+format, v...)
}

/*
Returns the path to this executable to call as the service.
*/
func exePath() string {
	exe, err := os.Executable()

	if err != nil {
		panic("Could not get executable path!")
	}

	return exe
}

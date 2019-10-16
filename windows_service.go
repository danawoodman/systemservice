// +build windows

package systemservice

import (
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

var blankManager = &mgr.Mgr{}
var blankService = &mgr.Service{}

/*
openManger opens the Windows service manager and returns the manager.
It is just a lightweight wrapper around "mgr.Connect()"
*/
func openManager() (m *mgr.Mgr, err error) {
	m, err = mgr.Connect()

	if err != nil {
		logger.Log("open manager error: ", err)
		return blankManager, err
	}

	// defer m.Disconnect()

	return m, nil
}

/*
connectService connects to a Window service by name and
returns the service or an error
*/
func connectService(name string) (s *mgr.Service, err error) {
	m, err := openManager()

	if err != nil {
		return blankService, err
	}

	s, err = m.OpenService(name)

	if err != nil {
		e := err.Error()
		logger.Log("open manager error: ", e)

		if strings.Contains(e, "specified service does not exist") {
			return blankService, &ServiceDoesNotExistError{serviceName: name}
		}

		return blankService, err
	}

	// defer s.Close()

	return s, nil
}

/*
runScCommand makes calls to the sc.exe binary.

See this page for reference:
https://www.computerhope.com/sc-command.htm
*/
func runScCommand(args ...string) (out string, err error) {
	logger.Log("running command: sc ", strings.Join(args, " "))
	return runCommand("sc", args...)
}

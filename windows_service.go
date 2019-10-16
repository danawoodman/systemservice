// +build windows

package systemservice

import (
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

var blankManager = &mgr.Mgr{}
var blankService = &mgr.Service{}

/*
openManger opens the Windows service manager and returns the manager
*/
func openManager() (m *mgr.Mgr, err error) {
	m, err = mgr.Connect()

	if err != nil {
		return blankManager, err
	}

	defer m.Disconnect()

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
		if strings.Contains(err.Error(), "specified service does not exist") {
			return blankService, &ServiceDoesNotExistError{serviceName: name}
		}

		return blankService, err
	}

	defer s.Close()

	return s, nil
}

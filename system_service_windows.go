// +build windows

package systemservice

import (
	"syscall"
	"time"
	// "golang.org/x/sys/windows/svc"
	// "golang.org/x/sys/windows/svc/debug"
	// "golang.org/x/sys/windows/svc/eventlog"
	// "golang.org/x/sys/windows/svc/mgr"
)

/*
Install the system service. If start is passed, also starts
the service.
*/
func (s *SystemService) Install(start bool) error {
	beep()

	logger.Log("install service")

	manager, err := openManager()

	if err != nil {
		return err
	}

	logger.Logf("manager: %+v", manager)

	return nil
}

/*
Start the system service if it is installed
*/
func (s *SystemService) Start() error {
	beep()
	return nil
}

/*
Restart attempts to stop the service if running then starts it again
*/
func (s *SystemService) Restart() error {
	beep()
	return nil
}

/*
Stop stops the system service by unloading the unit file
*/
func (s *SystemService) Stop() error {
	beep()
	return nil
}

/*
Uninstall the system service by first stopping it then removing
the unit file.
*/
func (s *SystemService) Uninstall() error {
	beep()
	return nil
}

/*
Status returns whether or not the system service is running
*/
func (s *SystemService) Status() (status ServiceStatus, err error) {
	beep()

	logger.Log("connect to service")

	serv, err := connectService("foo")

	if err != nil {
		return ServiceStatus{}, err
	}

	logger.Logf("service: %+v", serv)

	status = ServiceStatus{}
	return status, nil
}

/*
Return whether or not the unit file eixts
*/
func (s *SystemService) Exists() bool {
	beep()
	return false
}

var beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")

func beep() {
	beepFunc.Call(0xffffffff)
	time.Sleep(time.Second)
}

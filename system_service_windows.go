// +build windows

package systemservice

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
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
	logger.Log("install service")

	name := s.Command.Name
	prog := s.Command.String()
	args := []string{
		"create",
		fmt.Sprintf("\"%s\"", name),
		"binPath=",
		fmt.Sprintf("\"%s\"", prog),
		// "start=",
		// "boot",
	}

	out, err := runScCommand(args...)

	if err != nil {
		if strings.Contains(err.Error(), "exit status 1073") {
			logger.Log("service already exists")
		} else {
			logger.Log("sc create output:\n", out)
			return err
		}
	}

	if start {
		if err := s.Start(); err != nil {
			return err
		}
	}

	// if strings.Contains(out, "SUCCESS") {
	// 	return nil
	// }

	beep()

	return nil
}

/*
Start the system service if it is installed
*/
func (s *SystemService) Start() error {
	_, err := runScCommand("start", fmt.Sprintf("\"%s\"", s.Command.Name))

	if err != nil {
		logger.Log("start service error: ", err)
		return err
	}

	return nil
}

/*
Restart attempts to stop the service if running then starts it again
*/
func (s *SystemService) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}

	if err := s.Start(); err != nil {
		return err
	}

	return nil
}

/*
Stop stops the system service by unloading the unit file
*/
func (s *SystemService) Stop() error {
	_, err := runScCommand("stop", fmt.Sprintf("\"%s\"", s.Command.Name))

	if err != nil {
		logger.Log("stop service error: ", err)

		if strings.Contains(err.Error(), "exit status 1062") {
			logger.Log("service already stopped")
		} else {
			return err
		}
	}

	return nil
}

/*
Uninstall the system service by first stopping it then removing
the unit file.
*/
func (s *SystemService) Uninstall() error {
	name := s.Command.Name

	err := s.Stop()

	if err != nil {
		return err
	}

	_, err = runScCommand("delete", fmt.Sprintf("\"%s\"", name))

	if err != nil {
		logger.Log("delete service error: ", err)
		return err
	}

	return nil
}

/*
Status returns whether or not the system service is running
*/
func (s *SystemService) Status() (status ServiceStatus, err error) {
	name := s.Command.Name

	logger.Log("getting service status")

	out, _ := runScCommand("queryex", name)

	pid := 0
	running := false

	// Parse the output looking
	lines := strings.Split(strings.TrimSpace(string(out)), "\r")
	for _, line := range lines {

		// Get PID from output
		if strings.Contains(line, "PID") {
			chunks := strings.Split(line, ":")
			if p := chunks[1]; p != "" {
				cleaned, _ := strconv.Atoi(strings.TrimSpace(p))
				if cleaned != 0 {
					pid = cleaned
				}
			}
		}

		// Get service running status from output
		if strings.Contains(line, "STATE") && strings.Contains(line, "RUNNING") {
			running = true
		}
	}

	return ServiceStatus{Running: running, PID: pid}, nil
}

/*
Return whether or not the unit file eixts
*/
func (s *SystemService) Exists() bool {
	_, err := runScCommand("queryex", fmt.Sprintf("\"%s\"", s.Command.Name))

	if err != nil {
		logger.Log("exists service error: ", err)
		// Service does not exist
		// if strings.Contains(err.Error(), "FAILED 1060") {
		// return false
		// }
		return false
	}

	return true
}

var beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")

func beep() {
	beepFunc.Call(0xffffffff)
}

// logger.Log("uninstall service: ", name)

// serv, err := connectService(name)

// if err != nil {
// 	logger.Log("uninstall error: ", err)
// 	return err
// }

// defer serv.Close()

// logger.Log("deleting service")

// err = serv.Delete()

// if err != nil {
// 	logger.Log("delete service error: ", err)
// 	return err
// }

// logger.Log("service deleted")
// logger.Log("removing event log")

// err = eventlog.Remove(name)

// if err != nil {
// 	logger.Log("remove event log error: ", err)
// 	return err
// }

// logger.Log("event log removed")

// manager, err := openManager()

// if err != nil {
// 	return err
// }

// defer manager.Disconnect()

// cmd := s.Command

// serv, err := manager.CreateService(
// 	cmd.Name,
// 	cmd.Program,
// 	mgr.Config{DisplayName: cmd.Name},
// 	cmd.Args...,
// )

// if err != nil {
// 	return err
// }

// defer serv.Close()

// err = eventlog.InstallAsEventCreate(cmd.Name, eventlog.Error|eventlog.Warning|eventlog.Info)

// if err != nil {
// 	serv.Delete()
// 	return fmt.Errorf("SetupEventLogSource() failed: %s", err)
// }

// logger.Logf("manager: %+v", manager)

// serv, err := connectService(name)

// if err != nil {
// 	logger.Log("status error: ", err)
// 	return ServiceStatus{}, err
// }

// defer serv.Close()

// logger.Logf("service: %+v", serv)

// pid := getPID(name)

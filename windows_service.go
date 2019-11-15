// +build windows

package systemservice

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/mgr"
	// "golang.org/x/sys/windows/svc/eventlog"
)

/*
connectService connects to a Window service by name and
returns the service or an error
*/
func connectService(name string) (s *mgr.Service, err error) {
	m, err = mgr.Connect()

	if err != nil {
		logger.Log("open manager error: ", err)
		return nil, err
	}

	s, err = m.OpenService(name)

	if err != nil {
		e := err.Error()
		logger.Log("open manager error: ", e)

		if strings.Contains(e, "specified service does not exist") {
			return nil, &ServiceDoesNotExistError{serviceName: name}
		}

		return nil, err
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

var elog debug.Log

type windowsService struct{}

func (m *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	logger.Log("execute called")
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	// fasttick := time.Tick(500 * time.Millisecond)
	// slowtick := time.Tick(2 * time.Second)
	// tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		// logger.Log("Loop!")
		select {
		// case <-tick:
		// 	beep()
		// 	elog.Info(1, "beep")
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				testOutput := strings.Join(args, "-")
				testOutput += fmt.Sprintf("-%d", c.Context)
				elog.Info(1, testOutput)
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				// tick = slowtick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				// tick = fasttick
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

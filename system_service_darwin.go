// +build darwin

package systemservice

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Run is a no-op on Darwin based systems
*/
func (s *SystemService) Run() error {
	return nil
}

/*
Install the system service. If start is passed, also starts
the service.
*/
func (s *SystemService) Install(start bool) error {
	plist := newPlist(s)

	path := plist.Path()
	dir := filepath.Dir(path)

	logger.Log("making sure folder exists: ", dir)

	os.MkdirAll(dir, os.ModePerm)

	logger.Log("generating plist file")

	content, err := plist.Generate()

	if err != nil {
		return err
	}

	logger.Log("writing plist to: ", path)

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return err
	}

	logger.Log("wrote plist:\n", content)

	if start {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
Start the system service if it is installed
*/
func (s *SystemService) Start() error {
	plist := newPlist(s)

	logger.Log("loading plist with launchctl")

	_, err := runLaunchCtlCommand("load", "-w", plist.Path())

	if err != nil {
		e := strings.ToLower(err.Error())

		// If not installed, install the service and then run start again.
		if strings.Contains(e, "no such file or directory") {
			logger.Log("service not installed yet, installing...")

			err = s.Install(true)

			if err != nil {
				return err
			}
		}

		// We don't care if the process fails because it is already
		// loaded
		if strings.Contains(e, "service already loaded") {
			logger.Log("service already loaded")
			return nil
		}

		return err
	}

	return nil
}

/*
Restart attempts to stop the service if running then starts it again
*/
func (s *SystemService) Restart() error {
	err := s.Stop()

	if err != nil {
		return err
	}

	err = s.Start()

	if err != nil {
		return err
	}

	return nil
}

/*
Stop stops the system service by unloading the plist file
*/
func (s *SystemService) Stop() error {
	plist := newPlist(s)

	_, err := runLaunchCtlCommand("unload", "-w", plist.Path())

	if err != nil {
		e := strings.ToLower(err.Error())

		if strings.Contains(e, "could not find specified service") {
			logger.Log("no service matching plist running: ", plist.Label)
			return nil
		}

		if strings.Contains(e, "no such file or directory") {
			logger.Log("plist file doesn't exist, nothing to stop: ", plist.Label)
			return nil
		}

		return err
	}

	return nil
}

/*
Uninstall the system service by first stopping it then removing
the plist file.
*/
func (s *SystemService) Uninstall() error {
	err := s.Stop()

	if err != nil {
		// If there is no matching process, don't throw an error
		// as it is already stopped.
		if strings.Contains(err.Error(), "exit status 3") != true {
			return err
		}
	}

	plist := newPlist(s)

	logger.Log("remove plist file")

	err = os.Remove(plist.Path())

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no such file or directory") {
			return nil
		}

		return err
	}

	return nil
}

/*
Status returns whether or not the system service is running
*/
func (s *SystemService) Status() (status *ServiceStatus, err error) {
	plist := newPlist(s)

	list, err := runLaunchCtlCommand("list")

	status = &ServiceStatus{}

	if err != nil {
		logger.Log("error getting launchctl status: ", err)
		return status, err
	}

	lines := strings.Split(strings.TrimSpace(string(list)), "\n")
	pattern := plist.Label

	if pattern == "" {
		return status, err
	}

	// logger.Log("running services:")

	for _, line := range lines {

		// logger.Log("line: ", line)

		chunks := strings.Split(line, "\t")

		if chunks[2] == pattern {
			if chunks[0] != "-" {
				pid, err := strconv.Atoi(chunks[0])

				if err != nil {
					return status, err
				}
				status.PID = pid
			}

			if status.PID != 0 {
				status.Running = true
			}
		}
	}

	return status, nil
}

/*
Return whether or not the plist file eixts
*/
func (s *SystemService) Exists() bool {
	plist := newPlist(s)

	return fileExists(plist.Path())
}

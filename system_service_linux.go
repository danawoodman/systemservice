// +build linux

package systemservice

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Run is a no-op on Linux based systems
*/
func (s *SystemService) Run() error {
	return nil
}

/*
Install the system service. If start is passed, also starts
the service.
*/
func (s *SystemService) Install(start bool) error {
	unit := newUnitFile(s)

	path := unit.Path()
	dir := filepath.Dir(path)

	logger.Log("making sure folder exists: ", dir)

	os.MkdirAll(dir, os.ModePerm)

	logger.Log("generating unit file")

	content, err := unit.Generate()

	if err != nil {
		return err
	}

	logger.Log("writing unit to: ", path)

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return err
	}

	logger.Log("wrote unit:\n", content)

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
	unit := newUnitFile(s)

	logger.Log("loading unit file with systemd")

	_, err := runSystemCtlCommand("start", unit.Label)

	if err != nil {
		return err
	}

	logger.Log("enabling unit file with systemd")

	_, err = runSystemCtlCommand("enable", unit.Label)

	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Created symlink") {
			logger.Log("")
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
	unit := newUnitFile(s)

	_, err := runSystemCtlCommand("reload-or-restart", unit.Label)

	if err != nil {
		return err
	}

	return nil
}

/*
Stop stops the system service by unloading the unit file
*/
func (s *SystemService) Stop() error {
	unit := newUnitFile(s)

	logger.Log("stopping unit file with systemd")

	_, err := runSystemCtlCommand("stop --now", unit.Label)

	if err != nil {
		return err
	}

	logger.Log("disabling unit file with systemd")

	_, err = runSystemCtlCommand("disable", unit.Label)

	if err != nil {
		if strings.Contains(err.Error(), "Removed") {
			logger.Log("ignoring remove symlink error")
			return nil
		}
		return err
	}

	return nil
}

/*
Uninstall the system service by first stopping it then removing
the unit file.
*/
func (s *SystemService) Uninstall() error {
	err := s.Stop()

	if err != nil {
		return err
	}

	logger.Log("remove unit file")

	unit := newUnitFile(s)
	err = unit.Remove()

	if err != nil {
		return err
	}

	return nil
}

/*
Status returns whether or not the system service is running
*/
func (s *SystemService) Status() (status *ServiceStatus, err error) {
	unit := newUnitFile(s)
	active, _ := runSystemCtlCommand("is-active", unit.Label)

	status = &ServiceStatus{}

	// Check if service is running
	if !strings.Contains(active, "active") {
		return status, nil
	}

	stat, _ := runSystemCtlCommand("status", unit.Label)

	// Get the PID from the status output
	lines := strings.Split(stat, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Main PID") {
			parts := strings.Split(strings.TrimSpace(line), " ")
			pid, _ := strconv.Atoi(parts[2])
			if pid != 0 {
				status.PID = pid
			}
		}
	}

	status.Running = true

	return status, nil
}

/*
Exists returns whether or not the unit file eixts
*/
func (s *SystemService) Exists() bool {
	unit := newUnitFile(s)
	return fileExists(unit.Path())
}

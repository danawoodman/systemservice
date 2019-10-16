// +build linux

package systemservice

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Install the system service. If start is passed, also starts
the service.
*/
func (s *SystemService) Install(start bool) error {
	unit := newUnitFile(s)

	path := unit.Path()
	dir := filepath.Dir(path)

	log.Println("making sure folder exists: ", dir)

	os.MkdirAll(dir, os.ModePerm)

	log.Println("generating unit file")

	content, err := unit.Generate()

	if err != nil {
		return err
	}

	log.Println("writing unit to: ", path)

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		return err
	}

	log.Print("wrote unit:\n", content)

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

	log.Println("loading unit file with systemd")

	_, err := runSystemCtlCommand("start", unit.Label)

	if err != nil {
		return err
	}

	log.Println("enabling unit file with systemd")

	_, err = runSystemCtlCommand("enable", unit.Label)

	if err != nil {
		if strings.Contains(err.Error(), "Created symlink") {
			log.Println("")

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

	log.Println("stopping unit file with systemd")

	_, err := runSystemCtlCommand("stop", unit.Label)

	if err != nil {
		return err
	}

	log.Println("disabling unit file with systemd")

	_, err = runSystemCtlCommand("disable", unit.Label)

	if err != nil {
		if strings.Contains(err.Error(), "Removed") {
			log.Println("ignoring remove symlink error")
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

	log.Println("remove unit file")

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
func (s *SystemService) Status() (status ServiceStatus, err error) {
	unit := newUnitFile(s)
	active, _ := runSystemCtlCommand("is-active", unit.Label)

	status = ServiceStatus{}

	// Check if service is running
	if strings.Contains(active, "active") != true {
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
Return whether or not the unit file eixts
*/
func (s *SystemService) Exists() bool {
	unit := newUnitFile(s)
	return fileExists(unit.Path())
}

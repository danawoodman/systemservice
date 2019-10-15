package systemservice

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

/*
New creates a new system service manager instance.
*/
func New(cmd ServiceCommand) SystemService {
	serv := SystemService{Command: cmd}
	return serv
}

/*
SystemService represents a generic system service configuration
*/
type SystemService struct {
	Command ServiceCommand
}

/*
ServiceCommand represents the command the system service should run
*/
type ServiceCommand struct {
	// The human-friendly name of your service. Note: best to not include
	// spaces in the name.
	Name string

	// The label to use to identify the service. This must be unique
	// and should not include spaces.
	Label string

	// The name of the program to run
	Program string

	// The arguments to pass to the command. Optional.
	Args []string

	// The description of your service. Optional.
	Description string
}

func (c *ServiceCommand) String() string {
	s := c.Program
	if len(c.Args) > 0 {
		s = s + " " + strings.Join(c.Args, " ")
	}
	return s
}

/*
ServiceStatus is a generic representation of the service running on the system
*/
type ServiceStatus struct {
	Running bool
	PID     int
}

/*
RunCommand is a lightweight wrapper around exec.Command
*/
func RunCommand(name string, args ...string) error {
	// cmdString := name + " " + strings.Join(args, " ")

	// logger.Debug("[system_service] running command: ", cmdString)

	stderr := bytes.NewBuffer(make([]byte, 0))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = stderr

	// err := cmd.Start()
	err := cmd.Run()

	if err != nil {
		// logger.Debugf("[sytstem_service] error running command '%s': %s", cmdString, err)
		return err
	}

	// If stderr returns anything, record that as an error.
	if stderr.Len() != 0 {
		return errors.New(stderr.String())
	}

	// logger.Debug("[system_service] command succeeded: ", cmdString)

	return nil
}

/*
IsRoot returns whether or not the program was run as root

Always returns false on Windows because there is no
good way to detect root on Windows.
*/
func isRoot() bool {
	u, err := user.Current()

	if err != nil {
		return false
	}

	// On unix systems, root user either has the UID 0,
	// the GID 0 or both.
	return u.Uid == "0" || u.Gid == "0"
}

/*
HomeDir returns the home directory of the user or "/" if
we cannot determine it.
*/
func homeDir() string {
	u, err := user.Current()

	if err != nil {
		fmt.Println("User does not have a home directory!")
		return "/"
	}

	return u.HomeDir
}

/*
fileExists is a helper to return whether or not a give
file exists
*/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

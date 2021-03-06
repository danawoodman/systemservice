package systemservice

import (
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

	// The URL to your service documentation. Optional.
	Documentation string

	// Whether or not to turn on debug behavior
	Debug bool
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
Running indicates if the service is active and running
*/
func (s *SystemService) Running() (bool, error) {
	status, err := s.Status()

	if err != nil {
		return false, err
	}

	return status.Running, nil
}

/*
runCommand is a lightweight wrapper around exec.Command
*/
func runCommand(name string, args ...string) (out string, err error) {
	stdout, err := exec.Command(name, args...).Output()
	return string(stdout), err
}

/*
isRoot returns whether or not the program was run as root

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
homeDir returns the home directory of the user or "/" if
we cannot determine it.
*/
func homeDir() string {
	u, err := user.Current()

	if err != nil {
		logger.Log("User does not have a home directory!")
		return "/"
	}

	return u.HomeDir
}

/*
username returns the username of the current user
*/
func username() string {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	return user.Username
}

/*
fileExists is a helper to return whether or not a give
file exists
*/
func fileExists(filename string) bool {
	info, err := appFS.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

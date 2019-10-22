// +build darwin

package systemservice

import "strings"

func runLaunchCtlCommand(cmd string, args ...string) (out string, err error) {
	logger.Log("running command: launchctl ", strings.Join(args, " "))
	return runCommand("launchctl", args...)
}

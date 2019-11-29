// +build linux

package systemservice

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func runSystemCtlCommand(cmd string, label string) (out string, err error) {
	args := strings.Split(cmd, " ")

	if !isRoot() {
		args = append(args, "--user")
	}

	args = append(args, label)

	logger.Log("running command: systemctl", strings.Join(args, " "))

	return runCommand("systemctl", args...)
}

/*
unitFile represents a launchctl unitFile file
*/
type unitFile struct {
	Label         string
	Command       string
	Description   string
	Documentation string
	StdOutPath    string
	StdErrPath    string
	User          string
}

func newUnitFile(serv *SystemService) unitFile {
	cmd := serv.Command
	label := cmd.Label

	user := username()
	if isRoot() {
		user = "root"
	}

	unit := unitFile{
		Label:         label,
		Command:       cmd.String(),
		Description:   cmd.Description,
		Documentation: cmd.Documentation,
		User:          user,
	}

	return unit
}

func (u *unitFile) Generate() (string, error) {
	var tmpl bytes.Buffer
	t := template.Must(template.New("unitFile").Parse(unitFileTemplate()))
	if err := t.Execute(&tmpl, u); err != nil {
		return "", err
	}

	return tmpl.String(), nil
}

func (u *unitFile) Path() string {
	file := u.Label + ".service"

	if isRoot() {
		return filepath.Join("/etc/systemd/system", file)
	}

	return filepath.Join(homeDir(), ".config/systemd/user", file)
}

func (u *unitFile) Remove() error {
	return os.Remove(u.Path())
}

/*
unitFileTemplate generates the contents of the unitFile file.
*/
func unitFileTemplate() string {
	return `[Unit]
After=network.target
Description={{ .Description }}
Documentation={{ .Documentation }}

[Service]
ExecStart={{ .Command }}
Restart=on-failure
Type=simple

[Install]
WantedBy=multi-user.target
`
}

// StandardOutput=file:{{ .LogPath }}

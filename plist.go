// +build darwin

package systemservice

import (
	"bytes"
	"encoding/xml"
	"path/filepath"
	"text/template"
)

/*
plist represents a launchctl plist file
*/
type plist struct {
	Label            string
	Program          string
	ProgramArguments []string
	KeepAlive        bool
	RunAtLoad        bool
	StdOutPath       string
	StdErrPath       string
}

func newPlist(serv *SystemService) plist {
	label := serv.Command.Label
	name := serv.Command.Name
	logDir := filepath.Join(homeDir(), "Library/Logs", name)
	args := []string{serv.Command.Program}
	if len(serv.Command.Args) != 0 {
		args = append(args, serv.Command.Args...)
	}

	pl := plist{
		Label:            label,
		ProgramArguments: args,
		KeepAlive:        true,
		RunAtLoad:        true,
		StdOutPath:       filepath.Join(logDir, name+".stdout.log"),
		StdErrPath:       filepath.Join(logDir, name+".stderr.log"),
	}

	return pl
}

func (p *plist) Generate() (string, error) {
	var tmpl bytes.Buffer
	t := template.Must(template.New("launchdConfig").Parse(plistTemplate()))
	if err := t.Execute(&tmpl, p); err != nil {
		return "", err
	}

	return tmpl.String(), nil
}

func (p *plist) Path() string {
	label := p.Label + ".plist"
	if isRoot() {
		return filepath.Join("/Library/LaunchDaemons/", label)
	}

	return filepath.Join(homeDir(), "Library/LaunchAgents/", label)
}

func (p *plist) String() string {
	encoded, _ := xml.MarshalIndent(p, "", "  ")
	return string(encoded)
}

/*
plistTemplate generates the contents of the plist file.
*/
func plistTemplate() string {
	return `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
<plist version='1.0'>
  <dict>
    <key>Label</key><string>{{ .Label }}</string>{{ if .Program }}
    <key>Program</key><string>{{ .Program }}</string>{{ end }}
    {{ if .ProgramArguments }}<key>ProgramArguments</key>
    <array>{{ range $arg := .ProgramArguments }}
      <string>{{ $arg }}</string>{{ end }}
    </array>{{ end }}
    <key>StandardOutPath</key>
    <string>{{ .StdOutPath }}</string>
    <key>StandardErrorPath</key>
    <string>{{ .StdErrPath }}</string>
    <key>KeepAlive</key> <{{ .KeepAlive }}/>
    <key>RunAtLoad</key> <{{ .RunAtLoad }}/>
  </dict>
</plist>
`
}

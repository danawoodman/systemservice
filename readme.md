# systemservice

[![GoDoc](https://godoc.org/github.com/danawoodman/systemservice?status.svg)](https://godoc.org/github.com/danawoodman/systemservice) [![MIT License](http://img.shields.io/:license-mit-blue.svg)](https://github.com/danawoodman/systemservice/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/danawoodman/systemservice)](https://goreportcard.com/report/github.com/danawoodman/systemservice) [![Test Status](https://github.com/danawoodman/systemservice/workflows/Test/badge.svg)](https://github.com/danawoodman/systemservice/actions) ![code size](https://img.shields.io/github/languages/code-size/danawoodman/systemservice?style=flat-square)

**THIS IS A WORK IN PROGRESS AND IS NOT CURRENLTY USABLE**

> A cross-platform system service manager written in Go

Operating system support:

- Windows (via `sc.exe`)
- Mac (via `launchctl`)
- Linux: (via `systemd`)

## Install

```shell
go get -u github.com/danawoodman/systemservice
```

## Usage

Please see the [examples folder](/examples) for examples.

First, setup a new service:

```go
// Create a command to run.
cmd := systemservice.ServiceCommand{
  Label: "some-unique-id",
  Program: "echo",
  Args: []string{"Hello", "World", "!"},
}

// Create the service
serv := systemservice.New(cmd)

```

Now you can manage your service as needed:

```go
serv.Install(start bool) error
serv.Start() error
serv.Restart() error
serv.Stop() error
serv.Uninstall() error
serv.Status() (systemservice.ServiceStatus, error)
serv.Running() bool
```

These commands are the same no matter the operating system target.

### Platform Notes

#### Mac OSX (aka Darwin)

Replace `<LABEL>` and `<NAME>` with the values you setup in your `Command`.

- If running as a root user:
  - Service plist is located at `/Library/LaunchDaemons/<LABEL>.plist`
  - Stdout logs are sent to `/Library/Logs/<NAME>/<NAME>.stdout.log`
  - Stderr logs are send to `/Library/Logs/<NAME>/<NAME>.stderr.log`
- If running as a non-root user:
  - Service plist is located at `~/Library/LaunchAgents/<LABEL>.plist`
  - Stdout logs are sent to `~/Library/Logs/<NAME>/<NAME>.stdout.log`
  - Stderr logs are send to `~/Library/Logs/<NAME>/<NAME>.stderr.log`

#### Linux (Systemd)

- View logs with `journalctl -u <LABEL>`

## Similar project

- <https://github.com/kardianos/service>
- <https://github.com/jvehent/service-go>

## License

MIT. See [license.md](/license.md)

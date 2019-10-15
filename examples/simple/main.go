package main

import (
	"log"

	"github.com/danawoodman/systemservice"
)

var serv systemservice.SystemService

func main() {
	cmd := systemservice.ServiceCommand{
		Name:        "MyService",
		Label:       "com.myservice",
		Program:     "echo", // "echo"
		Args:        []string{"Hello!"},
		Description: "My systemservice test!",
	}

	log.Println("created command: ", cmd.String())

	serv = systemservice.New(cmd)

	logStatus()

	if serv.Exists() {
		log.Println("service exists, uninstalling...")

		if err := serv.Uninstall(); err != nil {
			panic(err)
		}

		log.Println("service uninstalled")
	}

	log.Println("running command: ", serv.Command.String())

	logStatus()

	log.Println("installing service")

	if err := serv.Install(true); err != nil {
		panic(err)
	}

	logStatus()

	log.Println("stopping service")

	if err := serv.Stop(); err != nil {
		panic(err)
	}

	logStatus()

	log.Println("starting service")

	if err := serv.Start(); err != nil {
		panic(err)
	}

	logStatus()

	log.Println("restarting service")

	if err := serv.Restart(); err != nil {
		panic(err)
	}

	logStatus()

	// log.Println("uninstall service")

	// if err := serv.Uninstall(); err != nil {
	// 	panic(err)
	// }

	// logStatus()
}

func logStatus() {
	stat, err := serv.Status()

	if err != nil {
		panic(err)
	}

	log.Printf("[STATUS] running: %t, pid: %d\n", stat.Running, stat.PID)
}

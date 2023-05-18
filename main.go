package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/sarkk0x0/orbiter/cmd"
	"github.com/sarkk0x0/orbiter/utils/hook"
	"github.com/sirupsen/logrus"
)

func main() {
	eventChannel := make(chan *logrus.Entry)
	logrus.AddHook(hook.NewChannelHook(eventChannel))

	c := cli.NewCLI("orbiter", "0.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"daemon": func() (cli.Command, error) {
			return &cmd.DaemonCmd{
				EventChannel: eventChannel,
			}, nil
		},
		"autoscaler ls": func() (cli.Command, error) {
			return &cmd.AutoscalerListCmd{}, nil
		},
		"system events": func() (cli.Command, error) {
			return &cmd.SystemEventsCmd{}, nil
		},
	}

	exitStatus, _ := c.Run()
	os.Exit(exitStatus)
}

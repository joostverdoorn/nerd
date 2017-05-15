package command

import (
	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"github.com/nerdalize/nerd/nerd/conf"
	"github.com/pkg/errors"
)

//QueueCreate command
type QueueCreate struct {
	*command
}

//QueueCreateFactory returns a factory method for the join command
func QueueCreateFactory() (cli.Command, error) {
	comm, err := newCommand("nerd queue create", "initialize a new queue for workers to consume tasks fro", "", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create command")
	}
	cmd := &QueueCreate{
		command: comm,
	}
	cmd.runFunc = cmd.DoRun

	return cmd, nil
}

//DoRun is called by run and allows an error to be returned
func (cmd *QueueCreate) DoRun(args []string) (err error) {
	config, err := conf.Read()
	if err != nil {
		HandleError(err)
	}

	bclient, err := NewClient(cmd.ui)
	if err != nil {
		HandleError(err)
	}

	out, err := bclient.CreateQueue(config.CurrentProject.Name)
	if err != nil {
		HandleError(err)
	}

	logrus.Infof("Queue Creation: %v", out)
	return nil
}

package sdk

import (
	"encoding/json"
)

type CommandRunnerInterface interface {
	RunCommand(command Command) (string, error)
}

type CommandRunner struct {
	client clientPointer
	lib    bitwardenLibrary
}

func NewCommandRunner(client clientPointer, lib bitwardenLibrary) *CommandRunner {
	return &CommandRunner{
		client: client,
		lib:    lib,
	}
}

func (c *CommandRunner) RunCommand(command Command) (string, error) {
	commandJSON, err := json.Marshal(command)
	if err != nil {
		return "", err
	}

	responseStr, err := c.lib.runCommand(string(commandJSON), c.client)
	if err != nil {
		return "", err
	}

	return responseStr, nil
}

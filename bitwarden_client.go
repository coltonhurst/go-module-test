package sdk

import (
	"encoding/json"
)

type BitwardenClient struct {
	client        clientPointer
	lib           bitwardenLibrary
	commandRunner CommandRunnerInterface
	Projects      ProjectsInterface
	Secrets       SecretsInterface
}

func NewBitwardenClient(apiURL *string, identityURL *string) (*BitwardenClient, error) {
	deviceType := DeviceType("SDK")
	userAgent := "Bitwarden GOLANG-SDK"
	clientSettings := ClientSettings{
		APIURL:      apiURL,
		IdentityURL: identityURL,
		UserAgent:   &userAgent,
		DeviceType:  &deviceType,
	}

	settingsJSON, err := json.Marshal(clientSettings)
	if err != nil {
		return nil, err
	}

	lib := newBitwardenLibrary()
	client, err := lib.init(string(settingsJSON))
	if err != nil {
		return nil, err
	}
	runner := NewCommandRunner(client, lib)

	return &BitwardenClient{
		lib:           lib,
		client:        client,
		commandRunner: runner,
		Projects:      NewProjects(runner),
		Secrets:       NewSecrets(runner),
	}, nil
}

func (c *BitwardenClient) AccessTokenLogin(accessToken string, statePath *string) error {
	req := AccessTokenLoginRequest{AccessToken: accessToken, StateFile: statePath}
	command := Command{AccessTokenLogin: &req}

	responseStr, err := c.commandRunner.RunCommand(command)
	if err != nil {
		return err
	}

	var response APIKeyLoginResponse
	return checkSuccessAndError(responseStr, &response)
}

func (c *BitwardenClient) Close() {
	c.lib.freeMem(c.client)
}

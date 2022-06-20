package common

import "encoding/base64"

type Base64EncodedValue string

func (field *Base64EncodedValue) SetValue(encoded string) error {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	*field = Base64EncodedValue(decoded)
	return nil
}

type github struct {
	AppID      string             `env:"GITHUB_APP_ID"`
	ClinetID   string             `env:"GITHUB_CLIENT_ID"`
	PrivateKey Base64EncodedValue `env:"GITHUB_PRIVATE_KEY"`
}

type conf struct {
	Github github
}

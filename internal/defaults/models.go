package defaults

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

type githubConf struct {
	AppID         string `env:"GITHUB_APP_ID"`
	ClinetID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
	// PrivateKey Base64EncodedValue `env:"GITHUB_PRIVATE_KEY"`
}

type kratosConf struct {
	ServerUrl string `env:"KRATOS_SERVER_URL"`
}

type conf struct {
	Github githubConf
	Kratos kratosConf
}

package conf

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

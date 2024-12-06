package cast

import (
	"encoding/json"

	"go.breu.io/quantm/internal/hooks/slack/defs"
)

func ByteToMessageProviderSlackUserInfo(data []byte) (*defs.MessageProviderSlackUserInfo, error) {
	d := &defs.MessageProviderSlackUserInfo{}

	err := json.Unmarshal(data, d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func ByteToMessageProviderSlackData(data []byte) (*defs.MessageProviderSlackData, error) {
	d := &defs.MessageProviderSlackData{}

	err := json.Unmarshal(data, d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

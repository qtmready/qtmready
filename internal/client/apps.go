// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package client

import (
	"go.breu.io/ctrlplane/internal/entities"
)

func (c *Client) AppList() ([]entities.App, error) {
	url := "/apps"
	reply := make([]entities.App, 0)

	if err := c.request("GET", url, &reply, nil); err != nil {
		return reply, err
	}

	return reply, nil
}

package models

import "github.com/scylladb/gocqlx/table"

var githubInstallMeta = table.Metadata{
	Name: "github_installations",
	Columns: []string{
		"id",
		"installation_id",
		"login",
		"type",
		"sender_id",
		"sender_login",
	},
}

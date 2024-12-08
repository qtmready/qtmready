package activities

import (
	"context"
	"fmt"
	"net/http"

	ghi "github.com/bradleyfalzon/ghinstallation/v2"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/config"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Kernel struct{}
)

func (k *Kernel) TokenizedCloneUrl(ctx context.Context, repo *entities.Repo) (string, error) {
	ghrepo, err := db.Queries().GetGithubRepoByID(ctx, repo.HookID)
	if err != nil {
		return "", err
	}

	install, err := db.Queries().GetGithubInstallation(ctx, ghrepo.InstallationID)
	if err != nil {
		return "", err
	}

	client, err := ghi.New(http.DefaultTransport, config.Instance().AppID, install.InstallationID, []byte(config.Instance().PrivateKey))
	if err != nil {
		return "", err
	}

	token, err := client.Token(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://git:%s@github.com/%s.git", token, ghrepo.FullName), nil
}

func (k *Kernel) DetectChanges(ctx context.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Push]) error {
	return nil
}

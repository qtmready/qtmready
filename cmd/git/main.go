package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/hooks/github"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/utils"
)

type (
	Config struct {
		Github *github.Config `koanf:"GITHUB"`
		DB     *db.Config     `koanf:"DB"`
	}
)

func main() {
	cfg := configure()
	ctx := context.Background()

	github.Configure(github.WithConfig(cfg.Github))
	kernel.Configure(
		kernel.WithRepoHook(eventsv1.RepoHook_REPO_HOOK_GITHUB, &github.KernelImpl{}),
	)

	db.Connection(db.WithConfig(cfg.DB))

	_ = db.Connection().Start(ctx)

	id := uuid.MustParse("019340e8-e115-7253-816b-2261d3128902")
	r, _ := db.Queries().GetRepo(ctx, id)

	path := utils.MustUUID().String()
	branch := "one"
	sha := "0c9b9b0aa97784a5cdfa2cc60d3e97d11def65ba"

	clone_pl := &defs.ClonePayload{Repo: &r, Hook: eventsv1.RepoHook_REPO_HOOK_GITHUB, Branch: branch, Path: path, SHA: sha}
	acts := repos.NewBranchActivities()
	path, _ = acts.Clone(ctx, clone_pl)
	rebased, _ := acts.Rebase(ctx, &defs.RebasePayload{Rebase: &eventsv1.Rebase{Base: branch, Head: sha}, Path: path})

	slog.Info("result", "result", rebased)

	_ = acts.RemoveDir(ctx, path)
}

func configure() *Config {
	config := &Config{}
	k := koanf.New("__")

	if err := k.Load(structs.Provider(config, "__"), nil); err != nil {
		panic(err)
	}

	// Load environment variables with the "__" delimiter.
	if err := k.Load(env.Provider("", "__", nil), nil); err != nil {
		panic(err)
	}

	// Unmarshal configuration from the Koanf instance to the Config struct.
	if err := k.Unmarshal("", config); err != nil {
		panic(err)
	}

	return config
}

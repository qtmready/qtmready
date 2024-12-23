package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/repos/git"
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
	sha := "0c9b9b0aa97784a5cdfa2cc60d3e97d11def65ba"

	r, err := db.Queries().GetRepo(ctx, id)
	if err != nil {
		slog.Error("failed to get repo from db", "error", err)
		os.Exit(1)

		return
	}

	path := utils.MustUUID().String()
	branch := "one"

	repo := git.NewRepository(&r, branch, path)

	err = repo.Clone(ctx)
	if err != nil {
		_ = err.(*git.RepositoryError).ReportError()

		os.Exit(1)

		return
	}

	slog.Info("repo cloned successfully", "path", path)

	diff, err := repo.Diff(ctx, branch, sha)
	if err != nil {
		_ = err.(git.GitError).ReportError()

		os.Exit(1)
	}

	fmt.Println(diff.Patch)

	slog.Info("diff generated successfully", "from", branch, "to", sha, "lines", diff.Lines)

	// Cleanup
	err = os.RemoveAll(path)
	if err != nil {
		slog.Error("failed to remove cloned directory", "path", path, "error", err)
	}
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

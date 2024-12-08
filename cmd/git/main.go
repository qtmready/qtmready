package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
	git "github.com/jeffwelling/git2go/v37"
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

	slog.Info("repo", "repo", r)

	path := utils.MustUUID().String()
	branch := "one"
	sha := "0c9b9b0aa97784a5cdfa2cc60d3e97d11def65ba"

	clone_pl := &defs.ClonePayload{Repo: &r, Hook: eventsv1.RepoHook_REPO_HOOK_GITHUB, Branch: branch, Path: path, SHA: sha}
	acts := repos.NewBranchActivities()
	path, _ = acts.Clone(ctx, clone_pl)

	rebase(path, sha, branch)

	// _ = acts.RemoveDir(ctx, path)
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

func rebase(path, sha, name string) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	// Lookup the commit to rebase onto
	id, err := git.NewOid(sha)
	if err != nil {
		log.Fatalf("Failed to create OID from commit hash: %v", err)
	}

	// Get the annotated commit for the commit to rebase onto
	head, err := repo.LookupAnnotatedCommit(id)
	if err != nil {
		log.Fatalf("Failed to get annotated commit: %v", err)
	}

	// Lookup the branch to rebase
	branch, err := repo.LookupBranch(name, git.BranchLocal)
	if err != nil {
		log.Fatalf("Failed to lookup branch: %v", err)
	}

	// Get the annotated commit for the branch
	upstream, err := repo.AnnotatedCommitFromRef(branch.Reference)
	if err != nil {
		log.Fatalf("Failed to get annotated commit: %v", err)
	}

	// Perform the rebase
	opts, err := git.DefaultRebaseOptions()
	if err != nil {
		log.Fatalf("Failed to get default rebase options: %v", err)
	}

	analysis, _, err := repo.MergeAnalysis([]*git.AnnotatedCommit{head})
	if err != nil {
		log.Fatalf("Failed to get merge analysis: %v", err)
	}

	if analysis == git.MergeAnalysisUpToDate {
		return
	}

	slog.Info("analysis", "analysis", analysis)

	rebase, err := repo.InitRebase(upstream, head, nil, &opts)
	if err != nil {
		log.Fatalf("Failed to initialize rebase: %v", err)
	}

	slog.Info("rebase operations", "count", rebase.OperationCount())

	for {
		operation, err := rebase.Next()
		if err != nil {
			if git.IsErrorCode(err, git.ErrorCodeIterOver) {
				break
			}

			log.Fatalf("Failed to get next rebase operation: %v", err)
		}

		// Apply the rebase operation
		if operation.Type == git.RebaseOperationPick {
			commit, err := repo.LookupCommit(operation.Id)
			if err != nil {
				log.Fatalf("Failed to lookup commit: %v", err)
			}

			index, err := repo.Index()
			if err != nil {
				log.Fatalf("Failed to get index: %v", err)
			}

			err = repo.CheckoutIndex(index, &git.CheckoutOptions{Strategy: git.CheckoutForce})
			if err != nil {
				log.Fatalf("Failed to checkout index: %v", err)
			}

			err = rebase.Commit(operation.Id, commit.Author(), commit.Committer(), commit.Message())
			if err != nil {
				log.Fatalf("Failed to create commit: %v", err)
			}
		}
	}

	err = rebase.Finish()
	if err != nil {
		log.Fatalf("Failed to finish rebase: %v", err)
	}

	fmt.Println("Rebase completed successfully")
}

package gitsourcerepo

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitSourceRepo struct {
	localPath string
	repo      *git.Repository
}

func Instance(_ context.Context, localPath string) *GitSourceRepo {
	return &GitSourceRepo{localPath: localPath}
}

func (s *GitSourceRepo) Clone(ctx context.Context, gitRepoUrl string, branchName string) error {
	repo, err := git.PlainCloneContext(ctx, s.localPath, false, &git.CloneOptions{
		URL:           gitRepoUrl,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		SingleBranch:  true,
		Progress:      nil,
	})
	s.repo = repo
	return err
}

func (s *GitSourceRepo) Path() string {
	return s.localPath
}
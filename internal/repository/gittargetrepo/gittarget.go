package gittargetrepo

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

type GitTargetRepo struct{
	localPath string
	repo      *git.Repository
}

var PushFunc = func(repo *git.Repository) error {
	// TODO need auth in PushOptions (secret for operator?)
	// return repo.Push(&git.PushOptions{})
	return nil
}

func Instance(_ context.Context, localPath string) *GitTargetRepo {
	return &GitTargetRepo{localPath: localPath}
}

func (t *GitTargetRepo) Clone(ctx context.Context, gitRepoUrl string) error {
	repo, err := git.PlainCloneContext(ctx, t.localPath, false, &git.CloneOptions{
		URL:           gitRepoUrl,
		Progress:      nil,
	})
	t.repo = repo
	return err
}

func (t *GitTargetRepo) GetHashForRevision(ctx context.Context, branchOrTag string) *plumbing.Hash {
	hash, err := t.repo.ResolveRevision(plumbing.Revision(branchOrTag))
	if err != nil {
		// not found
		return nil
	}
	return hash
}

func (t *GitTargetRepo) Checkout(ctx context.Context, branch string) error {
	worktree, err := t.repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)})
	if err != nil {
		return err
	}

	return nil
}

func (t *GitTargetRepo) CreateBranch(ctx context.Context, shortBranchName string, hash *plumbing.Hash) error {
	refName := plumbing.ReferenceName("refs/heads/" + shortBranchName)
	ref := plumbing.NewHashReference(refName, *hash)
	return t.repo.Storer.SetReference(ref)
}

func (t *GitTargetRepo) CommitAndPush(ctx context.Context, name string, email string, message string) error {
	worktree, err := t.repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(".")
	if err != nil {
		return err
	}

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = PushFunc(t.repo)
	if err != nil {
		return err
	}

	return nil
}

func (t *GitTargetRepo) Path() string {
	return t.localPath
}

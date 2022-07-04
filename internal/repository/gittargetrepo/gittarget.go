package gittargetrepo

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"time"
)

type GitTargetRepo struct {
	localPath string
	repo      *git.Repository
	remote    *git.Remote
	pushFunc  func(auth transport.AuthMethod) error
}

// note: push is disabled by default until we enable it

func Instance(_ context.Context, localPath string) *GitTargetRepo {
	return &GitTargetRepo{
		localPath: localPath,
		pushFunc: func(_ transport.AuthMethod) error {
			return nil
		},
	}
}

func (t *GitTargetRepo) PrepareInit(ctx context.Context, gitRepoUrl string, gitBranch string) error {
	repo, err := git.PlainInit(t.localPath, false)
	t.repo = repo

	if gitBranch != "master" {
		h := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName(gitBranch))
		if err = t.repo.Storer.SetReference(h); err != nil {
			return err
		}
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{gitRepoUrl},
	})
	t.remote = remote

	return err
}

func (t *GitTargetRepo) Clone(ctx context.Context, gitRepoUrl string, auth transport.AuthMethod) error {
	repo, err := git.PlainCloneContext(ctx, t.localPath, false, &git.CloneOptions{
		Auth:     auth,
		URL:      gitRepoUrl,
		Progress: nil,
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

func (t *GitTargetRepo) CommitAndPush(ctx context.Context, name string, email string, message string, auth transport.AuthMethod) error {
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

	err = t.pushFunc(auth)
	if err != nil {
		return err
	}

	return nil
}

func (t *GitTargetRepo) EnablePush() {
	t.pushFunc = func(auth transport.AuthMethod) error {
		return t.repo.Push(&git.PushOptions{
			Auth: auth,
		})
	}
}

func (t *GitTargetRepo) Path() string {
	return t.localPath
}

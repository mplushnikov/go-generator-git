package gitsourcerepo

import (
	"context"
)

type GitSourceRepo struct{
	localPath string
}

func Instance(_ context.Context, localPath string) *GitSourceRepo {
	return &GitSourceRepo{localPath: localPath}
}

package gittargetrepo

import (
	"context"
)

type GitTargetRepo struct{
	localPath string
}

func Instance(_ context.Context, localPath string) *GitTargetRepo {
	return &GitTargetRepo{localPath: localPath}
}

package tmpdir

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type TmpDir struct {
	basePath string
	dirName  string
}

func Instance(_ context.Context, basePath string) *TmpDir {
	return &TmpDir{basePath: basePath, dirName: randomDirName()}
}

func (t *TmpDir) Create(_ context.Context) error {
	path := filepath.Join(t.basePath, t.dirName)
	return os.Mkdir(path, os.ModePerm)
}

func (t *TmpDir) DeleteRecursive(_ context.Context) error {
	path := filepath.Join(t.basePath, t.dirName)
	return os.RemoveAll(path)
}

func (t *TmpDir) Path(_ context.Context) string {
	return filepath.Join(t.basePath, t.dirName)
}

// internal helpers

func randomDirName() string {
	randomUuid, err := uuid.NewRandom()
	if err != nil {
		return randomDirNameFallback()
	}
	return randomUuid.String()
}

func randomDirNameFallback() string {
	return fmt.Sprintf("%s-%d",
		time.Now().Format("20060102-150405-999999999"),
		rand.Intn(65535))
}

package pipeline

import (
	"os"
	"path/filepath"

	"github.com/dpastoor/qvm/internal/config"
	"github.com/dpastoor/qvm/internal/gh"
	"github.com/dpastoor/qvm/internal/unarchive"
)

// DownloadReleaseVersion downloads the appropriate tarball for a given platform
// platform should be windows | linux | darwin
func DownloadReleaseVersion(release string, platform string, progress bool) (string, error) {
	client := gh.NewClient(os.Getenv("GITHUB_PAT"))
	dlPath, err := gh.DownloadReleaseAsset(client, release, platform, progress)
	if err != nil {
		return "", err
	}
	dirPath := filepath.Join(config.GetPathToVersionsDir(), release)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return "", err
	}
	dlFile, err := os.Open(dlPath)
	if err != nil {
		return "", err
	}
	err = unarchive.Unarchive(dlFile, dirPath)
	if err != nil {
		return "", err
	}
	return dirPath, err
}

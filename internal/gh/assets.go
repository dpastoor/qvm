package gh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	log "github.com/sirupsen/logrus"
)

func getOsAssetSuffix(os string) osAssetSuffix {
	switch os {
	case "linux":
		return linuxamd64
	case "darwin":
		return macos
	case "windows":
		return win
	default:
		return unknown
	}
}

type osAssetSuffix int64

const (
	unknown osAssetSuffix = iota
	linuxamd64
	macos
	win
)

func (o osAssetSuffix) String() string {
	switch o {
	case linuxamd64:
		return "linux-amd64.tar.gz"
	case macos:
		return "macos.tar.gz"
	case win:
		return "win.zip"
	default:
		return "unknown"
	}
}

// DownloadReleaseAsset downloads the release asset for a given platform to a temp
// file and returns the path to the written file.
// targetOs should be "windows", "darwin", "linux"
func DownloadReleaseAsset(client *github.Client, tag string, targetOs string) (string, error) {
	switch targetOs {
	case "windows", "darwin", "linux":
		break
	default:
		return "", fmt.Errorf("invalid target os: %s, must be one of linux,darwin,windows", targetOs)
	}
	release, err := GetRelease(client, tag)
	if err != nil {
		return "", err
	}
	asset := findAssetForOs(release.Assets, getOsAssetSuffix(targetOs))
	if asset == nil {
		return "", errors.New("no release asset found")
	}
	// shouldn't need the redirect url given should follow redirects with the http client
	log.Tracef("fetching information to download release asset from %s\n", asset.GetBrowserDownloadURL())
	start := time.Now()
	rc, _, err := client.Repositories.DownloadReleaseAsset(context.Background(), "quarto-dev", "quarto-cli", asset.GetID(), http.DefaultClient)
	log.Tracef("done fetching release asset information in %s\n", time.Since(start))
	if err != nil {
		return "", err
	}
	defer rc.Close()
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("*-%s", asset.GetName()))
	if err != nil {
		return tmpFile.Name(), err
	}
	defer tmpFile.Close()
	log.Tracef("starting to copy release asset to %s\n", tmpFile.Name())
	start = time.Now()
	wb, err := io.Copy(tmpFile, rc)
	log.Tracef("done copying release asset in %s\n", time.Since(start))
	if err != nil {
		return tmpFile.Name(), err
	}
	log.Debugf("downloaded %d bytes from release asset at url %s\n", wb, asset.GetBrowserDownloadURL())
	return tmpFile.Name(), nil
}

func findAssetForOs(assets []*github.ReleaseAsset, suffix osAssetSuffix) *github.ReleaseAsset {
	if suffix == unknown {
		return nil
	}
	for _, asset := range assets {
		name := asset.GetName()
		if strings.HasSuffix(name, suffix.String()) {
			return asset
		}
	}
	return nil
}

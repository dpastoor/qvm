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

	"github.com/dustin/go-humanize"
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
	case "rhel7":
		return rhel7
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
	rhel7
)

func (o osAssetSuffix) String() string {
	switch o {
	case linuxamd64:
		return "linux-amd64.tar.gz"
	case macos:
		return "macos.tar.gz"
	case win:
		return "win.zip"
	case rhel7:
		return "linux-rhel7-amd64.tar.gz"
	default:
		return "unknown"
	}
}

type WriteCounter struct {
	Written  int
	Total    int
	Label    string
	nWrites  int
	progress bool
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Written += n
	wc.nWrites++
	if wc.progress && wc.nWrites%300 == 0 {
		log.Infof("downloading quarto version %s ... (%s/%s)", wc.Label,
			humanize.Bytes(uint64(wc.Written)),
			humanize.Bytes(uint64(wc.Total)))
	}
	if wc.Written == wc.Total {
		log.Infof("completed downloading quarto version %s", wc.Label)
	}
	return n, nil
}

// DownloadReleaseAsset downloads the release asset for a given platform to a temp
// file and returns the path to the written file.
// targetOs should be "windows", "darwin", "linux"
func DownloadReleaseAsset(client *github.Client, tag string, targetOs string, progress bool) (string, error) {
	switch targetOs {
	case "windows", "darwin", "linux", "rhel7":
		break
	default:
		return "", fmt.Errorf("invalid target os: %s, must be one of linux,darwin,windows,rhel7", targetOs)
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
	counter := &WriteCounter{Total: asset.GetSize(), Label: tag, progress: progress}
	wb, err := io.Copy(tmpFile, io.TeeReader(rc, counter))
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

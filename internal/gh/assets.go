package gh

import (
	"strings"

	"github.com/google/go-github/v44/github"
)

type OsAssetSuffix int64

const (
	unknown OsAssetSuffix = iota
	linuxamd64
	macos
	win
)

func (o OsAssetSuffix) String() string {
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

func GetTarballForOs(assets []github.ReleaseAsset, suffix OsAssetSuffix) string {
	if suffix == unknown {
		return ""
	}
	var dlUrl string
	for _, asset := range assets {
		name := asset.GetName()
		if strings.HasSuffix(name, suffix.String()) {
			dlUrl = asset.GetBrowserDownloadURL()
			break
		}
	}
	return dlUrl
}

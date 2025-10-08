package gh

import (
	"runtime"
	"testing"

	"github.com/google/go-github/v73/github"
	"github.com/stretchr/testify/assert"
)

func TestGetOsAssetSuffix(t *testing.T) {
	// Note: Linux suffix depends on runtime.GOARCH, so we need to account for that
	expectedLinuxSuffix := "linux-amd64.tar.gz"
	if runtime.GOARCH == "arm64" {
		expectedLinuxSuffix = "linux-arm64.tar.gz"
	}

	tests := []struct {
		name     string
		os       string
		expected string
	}{
		{
			name:     "linux (arch-dependent)",
			os:       "linux",
			expected: expectedLinuxSuffix,
		},
		{
			name:     "darwin (macos)",
			os:       "darwin",
			expected: "macos.tar.gz",
		},
		{
			name:     "windows",
			os:       "windows",
			expected: "win.zip",
		},
		{
			name:     "rhel7",
			os:       "rhel7",
			expected: "linux-rhel7-amd64.tar.gz",
		},
		{
			name:     "unknown",
			os:       "freebsd",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suffix := getOsAssetSuffix(tt.os)
			assert.Equal(t, tt.expected, suffix.String())
		})
	}
}

func TestFindAssetForOs(t *testing.T) {
	// Create mock assets
	assets := []*github.ReleaseAsset{
		{Name: github.String("quarto-1.4.553-linux-amd64.tar.gz")},
		{Name: github.String("quarto-1.4.553-linux-arm64.tar.gz")},
		{Name: github.String("quarto-1.4.553-macos.tar.gz")},
		{Name: github.String("quarto-1.4.553-win.zip")},
		{Name: github.String("quarto-1.4.553-linux-rhel7-amd64.tar.gz")},
	}

	tests := []struct {
		name      string
		suffix    osAssetSuffix
		wantFound bool
		wantName  string
	}{
		{
			name:      "find linux amd64",
			suffix:    linuxamd64,
			wantFound: true,
			wantName:  "quarto-1.4.553-linux-amd64.tar.gz",
		},
		{
			name:      "find linux arm64",
			suffix:    linuxarm64,
			wantFound: true,
			wantName:  "quarto-1.4.553-linux-arm64.tar.gz",
		},
		{
			name:      "find macos",
			suffix:    macos,
			wantFound: true,
			wantName:  "quarto-1.4.553-macos.tar.gz",
		},
		{
			name:      "find windows",
			suffix:    win,
			wantFound: true,
			wantName:  "quarto-1.4.553-win.zip",
		},
		{
			name:      "find rhel7",
			suffix:    rhel7,
			wantFound: true,
			wantName:  "quarto-1.4.553-linux-rhel7-amd64.tar.gz",
		},
		{
			name:      "unknown suffix returns nil",
			suffix:    unknown,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := findAssetForOs(assets, tt.suffix)
			if tt.wantFound {
				assert.NotNil(t, asset)
				assert.Equal(t, tt.wantName, asset.GetName())
			} else {
				assert.Nil(t, asset)
			}
		})
	}
}

func TestFindAssetForOs_MissingAsset(t *testing.T) {
	// Simulate older release without ARM64 (like v1.2.335)
	assetsWithoutARM64 := []*github.ReleaseAsset{
		{Name: github.String("quarto-1.2.335-linux-amd64.tar.gz")},
		{Name: github.String("quarto-1.2.335-macos.tar.gz")},
		{Name: github.String("quarto-1.2.335-win.zip")},
	}

	asset := findAssetForOs(assetsWithoutARM64, linuxarm64)
	assert.Nil(t, asset, "Should return nil when ARM64 asset doesn't exist")
}

func TestGetOsAssetSuffix_LinuxArch(t *testing.T) {
	// This test verifies that the function correctly detects ARM64 vs AMD64
	// Note: This will match the actual runtime architecture where tests run
	suffix := getOsAssetSuffix("linux")

	if runtime.GOARCH == "arm64" {
		assert.Equal(t, "linux-arm64.tar.gz", suffix.String())
	} else {
		assert.Equal(t, "linux-amd64.tar.gz", suffix.String())
	}
}

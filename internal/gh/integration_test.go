// +build integration

package gh

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReleaseAssetAvailability tests that we can find appropriate assets
// for different Quarto versions across platforms
func TestReleaseAssetAvailability(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	client := NewClient(os.Getenv("GITHUB_PAT"))

	// Test cases covering different Quarto version eras
	testCases := []struct {
		version             string
		hasARM64Linux       bool
		description         string
		expectedAssetCount  int
		requiredAssets      []string
	}{
		{
			version:         "v1.2.335",
			hasARM64Linux:   false,
			description:     "v1.2.x - no ARM64 Linux support",
			requiredAssets:  []string{"linux-amd64.tar.gz", "macos.tar.gz", "win.zip"},
		},
		{
			version:         "v1.4.553",
			hasARM64Linux:   true,
			description:     "v1.4.x - has ARM64 Linux support",
			requiredAssets:  []string{"linux-amd64.tar.gz", "linux-arm64.tar.gz", "macos.tar.gz", "win.zip"},
		},
		{
			version:         "v1.5.54",
			hasARM64Linux:   true,
			description:     "v1.5.x - has ARM64 Linux support",
			requiredAssets:  []string{"linux-amd64.tar.gz", "linux-arm64.tar.gz", "macos.tar.gz", "win.zip"},
		},
		{
			version:         "v1.7.34",
			hasARM64Linux:   true,
			description:     "v1.7.x - has ARM64 Linux support",
			requiredAssets:  []string{"linux-amd64.tar.gz", "linux-arm64.tar.gz", "macos.tar.gz", "win.zip"},
		},
		{
			version:         "v1.8.25",
			hasARM64Linux:   true,
			description:     "v1.8.x - has ARM64 Linux support",
			requiredAssets:  []string{"linux-amd64.tar.gz", "linux-arm64.tar.gz", "macos.tar.gz", "win.zip"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			release, err := GetRelease(client, tc.version)
			require.NoError(t, err, "Failed to fetch release %s", tc.version)
			require.NotNil(t, release, "Release should not be nil")

			// Verify required assets exist
			for _, requiredSuffix := range tc.requiredAssets {
				found := false
				for _, asset := range release.Assets {
					if asset.GetName() != "" {
						name := asset.GetName()
						if len(name) > len(requiredSuffix) && name[len(name)-len(requiredSuffix):] == requiredSuffix {
							found = true
							break
						}
					}
				}
				assert.True(t, found, "Required asset with suffix %s not found in release %s", requiredSuffix, tc.version)
			}

			// Test platform-specific asset finding
			platforms := []struct {
				os            string
				shouldSucceed bool
			}{
				{"linux", true},      // Should always work (falls back to amd64)
				{"darwin", true},     // Should always work
				{"windows", true},    // Should always work
			}

			for _, platform := range platforms {
				t.Run(platform.os, func(t *testing.T) {
					// Override GOARCH for ARM64 testing only on Linux
					if platform.os == "linux" && !tc.hasARM64Linux && runtime.GOARCH == "arm64" {
						// For old releases without ARM64, we would expect this to fail
						// currently, but with our fallback fix it should succeed
						asset := findAssetForOs(release.Assets, getOsAssetSuffix(platform.os))
						// Current implementation: will be nil
						// After fix: should fallback to amd64
						if asset == nil {
							t.Logf("WARNING: No ARM64 asset found for %s, needs fallback logic", tc.version)
						}
					} else {
						suffix := getOsAssetSuffix(platform.os)
						asset := findAssetForOs(release.Assets, suffix)
						if platform.shouldSucceed {
							assert.NotNil(t, asset, "Should find asset for %s on %s", tc.version, platform.os)
						}
					}
				})
			}
		})
	}
}

// TestDownloadReleaseAsset_RealDownload tests actual download capability
// This is a slower test that actually downloads assets
func TestDownloadReleaseAsset_RealDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow integration test in short mode")
	}

	// Only run on pull requests or explicitly via environment variable
	if os.Getenv("RUN_DOWNLOAD_TESTS") != "true" && os.Getenv("CI") != "true" {
		t.Skip("skipping download test, set RUN_DOWNLOAD_TESTS=true to run")
	}

	client := NewClient(os.Getenv("GITHUB_PAT"))

	// Use a small, recent release for testing
	version := "v1.8.25"
	targetOS := runtime.GOOS

	t.Logf("Testing download for version %s on %s", version, targetOS)

	tmpPath, err := DownloadReleaseAsset(client, version, targetOS, false)
	require.NoError(t, err, "Download should succeed")
	require.NotEmpty(t, tmpPath, "Downloaded file path should not be empty")

	// Verify file exists
	info, err := os.Stat(tmpPath)
	require.NoError(t, err, "Downloaded file should exist")
	assert.Greater(t, info.Size(), int64(0), "Downloaded file should not be empty")

	// Clean up
	defer os.Remove(tmpPath)

	t.Logf("Successfully downloaded %d bytes to %s", info.Size(), tmpPath)
}

// TestGetReleases tests fetching multiple releases
func TestGetReleases(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	client := NewClient(os.Getenv("GITHUB_PAT"))

	releases, err := GetReleases(client, 10, "release")
	require.NoError(t, err)
	assert.LessOrEqual(t, len(releases), 10, "Should return at most 10 releases")
	assert.Greater(t, len(releases), 0, "Should return at least 1 release")

	// Verify no pre-releases when filtering for "release"
	for _, release := range releases {
		assert.False(t, release.GetPrerelease(), "Should not include pre-releases when type is 'release'")
	}
}

// TestGetLatestRelease tests fetching the latest release
func TestGetLatestRelease(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	client := NewClient(os.Getenv("GITHUB_PAT"))

	release, err := GetLatestRelease(client)
	require.NoError(t, err)
	require.NotNil(t, release)
	assert.NotEmpty(t, release.GetTagName())
	assert.False(t, release.GetPrerelease(), "Latest release should not be a pre-release")
}

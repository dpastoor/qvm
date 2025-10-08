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
		unsupportedOnARM64  bool // Should error on ARM64 Linux
	}{
		{
			version:            "v1.2.335",
			hasARM64Linux:      false,
			description:        "v1.2.x - no ARM64 Linux support (should error on ARM64)",
			requiredAssets:     []string{"linux-amd64.tar.gz", "macos.tar.gz", "win.zip"},
			unsupportedOnARM64: true,
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

			// Test platform-specific asset finding with strict behavior
			platforms := []struct {
				os            string
				shouldSucceed bool
			}{
				{"linux", true},      // Should work if asset exists for current arch
				{"darwin", true},     // Should always work
				{"windows", true},    // Should always work
			}

			for _, platform := range platforms {
				t.Run(platform.os, func(t *testing.T) {
					suffix := getOsAssetSuffix(platform.os)
					asset := findAssetForOs(release.Assets, suffix)

					// Special handling for ARM64 Linux on old releases
					if platform.os == "linux" && runtime.GOARCH == "arm64" && tc.unsupportedOnARM64 {
						// This should correctly fail - no ARM64 asset exists
						assert.Nil(t, asset, "Expected no ARM64 asset for %s (strict behavior)", tc.version)

						// Verify the error message would be helpful
						if asset == nil {
							t.Logf("✓ Correctly detected missing ARM64 asset for %s. Error would inform user.", tc.version)
						}
					} else {
						// All other cases should succeed
						assert.NotNil(t, asset, "Should find asset for %s on %s", tc.version, platform.os)
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

// TestDownloadReleaseAsset_UnsupportedPlatform tests strict error behavior
// for unsupported platform/architecture combinations
func TestDownloadReleaseAsset_UnsupportedPlatform(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Only test on ARM64 to verify strict behavior
	if runtime.GOARCH != "arm64" {
		t.Skip("This test only applies to ARM64 systems")
	}

	client := NewClient(os.Getenv("GITHUB_PAT"))

	// Try to download v1.2.335 on ARM64 Linux - should fail with helpful error
	tmpPath, err := DownloadReleaseAsset(client, "v1.2.335", "linux", false)

	// Should error
	assert.Error(t, err, "Should error when trying to download unsupported ARM64 release")
	assert.Empty(t, tmpPath, "Should not return a path when download fails")

	// Error should be informative
	if err != nil {
		errMsg := err.Error()
		assert.Contains(t, errMsg, "arm64", "Error should mention ARM64 architecture")
		assert.Contains(t, errMsg, "no release asset found", "Error should explain what went wrong")
		t.Logf("✓ Helpful error message: %s", errMsg)
	}
}

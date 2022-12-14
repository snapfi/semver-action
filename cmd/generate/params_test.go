package generate_test

import (
	"os"
	"testing"

	"github.com/snapfi/semver-action/cmd/generate"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadParams_Prefix(t *testing.T) {
	os.Setenv("INPUT_PREFIX", "ver")
	defer os.Unsetenv("INPUT_PREFIX")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "ver", params.Prefix)
}

func TestLoadParams_Prefix_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "v", params.Prefix)
}

func TestLoadParams_PrereleaseID(t *testing.T) {
	os.Setenv("INPUT_PRERELEASE_ID", "alpha")
	defer os.Unsetenv("INPUT_PRERELEASE_ID")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "alpha", params.PrereleaseID)
}

func TestLoadParams_PrereleaseID_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "pre", params.PrereleaseID)
}

func TestLoadParams_ForcePrerelease(t *testing.T) {
	os.Setenv("INPUT_FORCE_PRERELEASE", "true")
	defer os.Unsetenv("INPUT_FORCE_PRERELEASE")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.True(t, params.ForcePrerelease)
}

func TestLoadParams_ForcePrerelease_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.False(t, params.ForcePrerelease)
}

func TestLoadParams_BranchName(t *testing.T) {
	os.Setenv("INPUT_BRANCH_NAME", "master")
	defer os.Unsetenv("INPUT_BRANCH_NAME")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "master", params.BranchName)
}

func TestLoadParams_BranchName_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "main", params.BranchName)
}

func TestLoadParams_CommitSha(t *testing.T) {
	os.Setenv("GITHUB_SHA", "2f08f7b455ec64741d135216d19d7e0c4dd46458")
	defer os.Unsetenv("GITHUB_SHA")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "2f08f7b455ec64741d135216d19d7e0c4dd46458", params.CommitSha)
}

func TestLoadParams_InvalidCommitSha(t *testing.T) {
	os.Setenv("GITHUB_SHA", "any")
	defer os.Unsetenv("GITHUB_SHA")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_RepoDir(t *testing.T) {
	os.Setenv("INPUT_REPO_DIR", "/var/tmp/wakatime-cli")
	defer os.Unsetenv("INPUT_REPO_DIR")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, "/var/tmp/wakatime-cli", params.RepoDir)
}

func TestLoadParams_RepoDir_Default(t *testing.T) {
	params, err := generate.LoadParams()
	require.NoError(t, err)

	assert.Equal(t, ".", params.RepoDir)
}

func TestLoadParams_Bump(t *testing.T) {
	tests := map[string]string{
		"auto":  "auto",
		"major": "major",
		"minor": "minor",
		"patch": "patch",
		"empty": "auto",
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv("INPUT_BUMP", value)
			defer os.Unsetenv("INPUT_BUMP")

			params, err := generate.LoadParams()
			require.NoError(t, err)

			assert.Equal(t, value, params.Bump)
		})
	}
}

func TestLoadParams_InvalidBump(t *testing.T) {
	os.Setenv("INPUT_BUMP", "invalid")
	defer os.Unsetenv("INPUT_BUMP")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

func TestLoadParams_BaseVersion(t *testing.T) {
	os.Setenv("INPUT_BASE_VERSION", "1.2.3")
	defer os.Unsetenv("INPUT_BASE_VERSION")

	params, err := generate.LoadParams()
	require.NoError(t, err)

	var expected = semver.MustParse("1.2.3")

	assert.True(t, expected.EQ(*params.BaseVersion))
}

func TestLoadParams_InvalidBaseVersion(t *testing.T) {
	os.Setenv("INPUT_BASE_VERSION", "invalid")
	defer os.Unsetenv("INPUT_BASE_VERSION")

	_, err := generate.LoadParams()
	require.Error(t, err)
}

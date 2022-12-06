package generate_test

import (
	"errors"
	"testing"

	"github.com/snapfi/semver-action/cmd/generate"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	tests := map[string]struct {
		CurrentBranch string
		LatestTag     string
		AncestorTag   string
		SourceBranch  string
		Params        generate.Params
		Result        generate.Result
	}{
		"no previous tag": {
			CurrentBranch: "main",
			SourceBranch:  "major/some",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "auto",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v0.0.0",
				AncestorTag:  "",
				SemverTag:    "v1.0.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"doc branch into main": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1-alpha.1",
			AncestorTag:   "v0.2.0-alpha.1",
			SourceBranch:  "doc/some",
			Params: generate.Params{
				CommitSha:    "81918ffc",
				Bump:         "auto",
				Prefix:       "v",
				PrereleaseID: "alpha",
				BranchName:   "main",
			},
			Result: generate.Result{},
		},
		"feature branch into main": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1",
			SourceBranch:  "feature/some",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "auto",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1",
				SemverTag:    "v0.3.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"bugfix branch into main": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1",
			AncestorTag:   "v0.2.1-alpha.2",
			SourceBranch:  "bugfix/some",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "auto",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1",
				AncestorTag:  "v0.2.1-alpha.2",
				SemverTag:    "v0.2.2-alpha.1",
				IsPrerelease: true,
			},
		},
		"misc branch into main": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1-alpha.1",
			AncestorTag:   "v0.2.0-alpha.1",
			SourceBranch:  "misc/some",
			Params: generate.Params{
				CommitSha:    "81918ffc",
				Bump:         "auto",
				Prefix:       "v",
				PrereleaseID: "alpha",
				BranchName:   "main",
			},
			Result: generate.Result{},
		},
		"valid branch into main with with force_prerelease false and with pre-release latest tag": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1-alpha.1",
			AncestorTag:   "v0.2.0-alpha.1",
			SourceBranch:  "feature/some",
			Params: generate.Params{
				CommitSha:    "81918ffc",
				Bump:         "auto",
				Prefix:       "v",
				PrereleaseID: "alpha",
				BranchName:   "main",
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1-alpha.1",
				AncestorTag:  "v0.2.0-alpha.1",
				SemverTag:    "v0.3.0",
				IsPrerelease: false,
			},
		},
		"valid branch into main with with force_prerelease false": {
			CurrentBranch: "main",
			LatestTag:     "v0.2.1",
			AncestorTag:   "v0.2.0",
			SourceBranch:  "feature/some",
			Params: generate.Params{
				CommitSha:    "81918ffc",
				Bump:         "auto",
				Prefix:       "v",
				PrereleaseID: "alpha",
				BranchName:   "main",
			},
			Result: generate.Result{
				PreviousTag:  "v0.2.1",
				AncestorTag:  "v0.2.0",
				SemverTag:    "v0.3.0",
				IsPrerelease: false,
			},
		},
		"base version set": {
			CurrentBranch: "main",
			LatestTag:     "v2.6.19",
			SourceBranch:  "feature/semver-initial",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "auto",
				BaseVersion:     newSemVerPtr(t, "4.2.0"),
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19",
				SemverTag:    "v4.3.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"force bump major": {
			CurrentBranch: "main",
			LatestTag:     "v2.6.19-alpha.1",
			SourceBranch:  "semver-initial",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "major",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-alpha.1",
				SemverTag:    "v3.0.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"force bump minor": {
			CurrentBranch: "main",
			LatestTag:     "v2.6.19-alpha.1",
			SourceBranch:  "semver-initial",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "minor",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-alpha.1",
				SemverTag:    "v2.7.0-alpha.1",
				IsPrerelease: true,
			},
		},
		"force bump patch": {
			CurrentBranch: "main",
			LatestTag:     "v2.6.19-alpha.1",
			SourceBranch:  "semver-initial",
			Params: generate.Params{
				CommitSha:       "81918ffc",
				Bump:            "patch",
				Prefix:          "v",
				PrereleaseID:    "alpha",
				ForcePrerelease: true,
				BranchName:      "main",
			},
			Result: generate.Result{
				PreviousTag:  "v2.6.19-alpha.1",
				SemverTag:    "v2.6.20-alpha.1",
				IsPrerelease: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gc := initGitClientMock(
				t,
				test.LatestTag,
				test.AncestorTag,
				test.CurrentBranch,
				test.SourceBranch,
				test.Params.CommitSha,
			)

			result, err := generate.Tag(test.Params, gc)
			require.NoError(t, err)

			assert.Equal(t, test.Result, result)
		})
	}
}

func TestTag_InvalidBranchName(t *testing.T) {
	params := generate.Params{
		CommitSha:    "81918ffc",
		Bump:         "auto",
		Prefix:       "v",
		PrereleaseID: "alpha",
		BranchName:   "main",
	}

	gc := initGitClientMock(
		t,
		"v2.6.19-alpha.1",
		"",
		"main",
		"semver-initial",
		"81918ffc",
	)

	result, err := generate.Tag(params, gc)

	assert.EqualError(t, err, "failed to determine bump strategy: invalid bump strategy")

	assert.Empty(t, result)
}

func TestTag_IsNotRepo(t *testing.T) {
	gc := &gitClientMock{
		MakeSafeFn: func() error {
			return nil
		},
		IsRepoFn: func() bool {
			return false
		},
	}

	_, err := generate.Tag(generate.Params{}, gc)
	require.Error(t, err)

	assert.EqualError(t, err, "current folder is not a git repository")
}

func TestTag_MakeSafeErr(t *testing.T) {
	gc := &gitClientMock{
		MakeSafeFn: func() error {
			return errors.New("error")
		},
	}

	_, err := generate.Tag(generate.Params{}, gc)
	require.Error(t, err)

	assert.EqualError(t, err, "failed to make safe: error")
}

type gitClientMock struct {
	CurrentBranchFn        func() (string, error)
	CurrentBranchFnInvoked int
	IsRepoFn               func() bool
	IsRepoFnInvoked        int
	MakeSafeFn             func() error
	MakeSafeFnInvoked      int
	LatestTagFn            func() string
	LatestTagFnInvoked     int
	AncestorTagFn          func(include, exclude, branch string) string
	AncestorTagFnInvoked   int
	SourceBranchFn         func(commitHash string) (string, error)
	SourceBranchFnInvoked  int
}

func initGitClientMock(
	t *testing.T,
	latestTag,
	ancestorTag,
	currentBranch,
	sourceBranch,
	expectedCommitHash string) *gitClientMock {
	return &gitClientMock{
		CurrentBranchFn: func() (string, error) {
			return currentBranch, nil
		},
		IsRepoFn: func() bool {
			return true
		},
		MakeSafeFn: func() error {
			return nil
		},
		LatestTagFn: func() string {
			return latestTag
		},
		AncestorTagFn: func(include, exclude, branch string) string {
			return ancestorTag
		},
		SourceBranchFn: func(commitHash string) (string, error) {
			assert.Equal(t, expectedCommitHash, commitHash)
			return sourceBranch, nil
		},
	}
}

func (m *gitClientMock) CurrentBranch() (string, error) {
	m.CurrentBranchFnInvoked++
	return m.CurrentBranchFn()
}

func (m *gitClientMock) MakeSafe() error {
	m.MakeSafeFnInvoked++
	return m.MakeSafeFn()
}

func (m *gitClientMock) IsRepo() bool {
	m.IsRepoFnInvoked++
	return m.IsRepoFn()
}

func (m *gitClientMock) LatestTag() string {
	m.LatestTagFnInvoked++
	return m.LatestTagFn()
}

func (m *gitClientMock) AncestorTag(include, exclude, branch string) string {
	m.AncestorTagFnInvoked++
	return m.AncestorTagFn(include, exclude, branch)
}

func (m *gitClientMock) SourceBranch(commitHash string) (string, error) {
	m.SourceBranchFnInvoked++
	return m.SourceBranchFn(commitHash)
}

func newSemVerPtr(t *testing.T, s string) *semver.Version {
	version, err := semver.New(s)
	require.NoError(t, err)

	return version
}

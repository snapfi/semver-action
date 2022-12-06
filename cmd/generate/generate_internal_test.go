package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetermineBumpStrategy(t *testing.T) {
	tests := map[string]struct {
		SourceBranch    string
		DestBranch      string
		Bump            string
		ExpectedMethod  string
		ExpectedVersion string
	}{
		"source branch bugfix, dest branch main and auto bump": {
			SourceBranch:    "bugfix/some",
			DestBranch:      "main",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "patch",
		},
		"source branch doc, dest branch main and auto bump": {
			SourceBranch:    "doc/some",
			DestBranch:      "main",
			Bump:            "auto",
			ExpectedMethod:  "",
			ExpectedVersion: "",
		},
		"source branch feature, dest branch main and auto bump": {
			SourceBranch:    "feature/some",
			DestBranch:      "main",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "minor",
		},
		"source branch major, dest branch main and auto bump": {
			SourceBranch:    "major/some",
			DestBranch:      "main",
			Bump:            "auto",
			ExpectedMethod:  "build",
			ExpectedVersion: "major",
		},
		"source branch misc, dest branch main and auto bump": {
			SourceBranch:    "misc/some",
			DestBranch:      "main",
			Bump:            "auto",
			ExpectedMethod:  "",
			ExpectedVersion: "",
		},
		"patch bump": {
			Bump:           "patch",
			ExpectedMethod: "patch",
		},
		"minor bump": {
			Bump:           "minor",
			ExpectedMethod: "minor",
		},
		"major bump": {
			Bump:           "major",
			ExpectedMethod: "major",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			method, version, err := determineBumpStrategy(test.Bump, test.SourceBranch, test.DestBranch, "main")
			require.NoError(t, err)

			assert.Equal(t, test.ExpectedMethod, method)
			assert.Equal(t, test.ExpectedVersion, version)
		})
	}
}

func TestDetermineBumpStrategy_InvalidSource(t *testing.T) {
	_, _, err := determineBumpStrategy("auto", "some-branch", "main", "main")

	assert.EqualError(t, err, "invalid bump strategy")
}

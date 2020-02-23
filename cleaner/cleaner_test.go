package cleaner

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/stretchr/testify/assert"
)

func TestGetManifestWithRef(t *testing.T) {
	untaggedManifest := google.ManifestInfo{
		Created: time.Now().AddDate(0, -1, 0),
		Tags:    []string{},
	}

	taggedManifest := google.ManifestInfo{
		Size:    5456,
		Created: time.Now(),
		Tags:    []string{"master", "v0.0.0"},
	}

	manifests := map[string]google.ManifestInfo{
		"ref1": untaggedManifest,
		"ref2": taggedManifest,
	}

	expected := []Manifest{
		{
			ref:          "ref2",
			manifestInfo: taggedManifest,
		},
		{
			ref:          "ref1",
			manifestInfo: untaggedManifest,
		},
	}
	assert.Equal(t, expected, getManifestsWithRef(manifests))
}

func TestByCreationTimeDesc(t *testing.T) {
	oldest := Manifest{
		ref: "3",
		manifestInfo: google.ManifestInfo{
			Created: time.Now().AddDate(0, -2, 0),
		},
	}
	middle := Manifest{
		ref: "1",
		manifestInfo: google.ManifestInfo{
			Created: time.Now().AddDate(0, -1, 0),
		},
	}
	newest := Manifest{
		ref: "2",
		manifestInfo: google.ManifestInfo{
			Created: time.Now(),
		},
	}
	manifests := []Manifest{middle, newest, oldest}
	expected := []Manifest{newest, middle, oldest}

	assert.Equal(t, expected, byCreationTimeDesc(manifests))
}

func TestShouldByDefaultDeleteUntagged(t *testing.T) {
	repository := Repository{}
	untagged := Manifest{
		ref:          "1",
		manifestInfo: google.ManifestInfo{},
	}
	delete := shouldDelete(0, untagged, repository, []*regexp.Regexp{})
	assert.True(t, delete)
}

func TestShouldByDefaultKeepTagged(t *testing.T) {
	repository := Repository{}
	tagged := Manifest{
		ref: "1",
		manifestInfo: google.ManifestInfo{
			Tags: []string{"tag"},
		},
	}
	delete := shouldDelete(0, tagged, repository, []*regexp.Regexp{})
	assert.False(t, delete)
}

func TestShouldKeepGracePeriod(t *testing.T) {
	repository := Repository{
		GracePeriod: 5,
	}
	inGrace := Manifest{
		ref: "1",
		manifestInfo: google.ManifestInfo{
			Created: time.Now().AddDate(0, 0, -2),
		},
	}
	outGrace := Manifest{
		ref: "2",
		manifestInfo: google.ManifestInfo{
			Created: time.Now().AddDate(0, 0, -10),
		},
	}

	assert.False(t, shouldDelete(0, inGrace, repository, []*regexp.Regexp{}))
	assert.True(t, shouldDelete(0, outGrace, repository, []*regexp.Regexp{}))
}

func TestKeepLastImages(t *testing.T) {
	repository := Repository{
		NumberToKeep: 5,
	}
	manifest := Manifest{}
	assert.False(t, shouldDelete(0, manifest, repository, []*regexp.Regexp{}))
	assert.True(t, shouldDelete(5, manifest, repository, []*regexp.Regexp{}))
	assert.True(t, shouldDelete(10, manifest, repository, []*regexp.Regexp{}))
}

func TestShouldDeleteTagged(t *testing.T) {
	repository := Repository{
		CleanTags: true,
	}
	tagged := Manifest{
		manifestInfo: google.ManifestInfo{
			Tags: []string{"tagged"},
		},
	}

	assert.True(t, shouldDelete(0, tagged, repository, []*regexp.Regexp{}))
}

func TestKeepTagsPatterns(t *testing.T) {
	repository := Repository{
		CleanTags: true,
		KeepTags: []string{
			"master", "latest", "^v(0|\\d+).(0|\\d+).(0|\\d+)$",
		},
	}
	patterns, _ := getPatterns(repository.KeepTags)

	matched := []Manifest{
		{
			manifestInfo: google.ManifestInfo{
				Tags: []string{"master", "latest"},
			},
		},
		{
			manifestInfo: google.ManifestInfo{
				Tags: []string{"v0.0.0"},
			},
		},
		{
			manifestInfo: google.ManifestInfo{
				Tags: []string{"v1.23.255"},
			},
		},
	}

	mismatched := []Manifest{
		{
			manifestInfo: google.ManifestInfo{
				Tags: []string{"non-matched"},
			},
		},
		{
			manifestInfo: google.ManifestInfo{
				Tags: []string{"v1.23.255-dirty"},
			},
		},
	}

	for _, manifest := range matched {
		assert.False(t, shouldDelete(0, manifest, repository, patterns))
	}

	for _, manifest := range mismatched {
		assert.True(t, shouldDelete(0, manifest, repository, patterns))
	}
}

func TestGetPatterns(t *testing.T) {
	expressions := []string{
		"master", "latest", "^v(0|\\d+).(0|\\d+).(0|\\d+)$",
	}
	_, err := getPatterns(expressions)
	assert.NoError(t, err)
}

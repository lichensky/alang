package cleaner

import (
	"regexp"
	"sort"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/sirupsen/logrus"
)

// Manifest with reference
type Manifest struct {
	manifestInfo google.ManifestInfo
	ref          string
}

// Clean images from GCR
func Clean(repository Repository, dry bool) error {
	repo, err := name.NewRepository(repository.Name)
	if err != nil {
		return err
	}

	auth, err := google.NewEnvAuthenticator()
	if err != nil {
		return err
	}

	tags, err := google.List(repo, google.WithAuth(auth))
	if err != nil {
		return err
	}

	tagsToKeepPatterns, err := getPatterns(repository.KeepTags)
	if err != nil {
		return err
	}

	manifests := getManifestsWithRef(tags.Manifests)
	for i, manifest := range manifests {
		if shouldDelete(i, manifest, repository, tagsToKeepPatterns) {
			deleteImage(manifest, dry)
		}
	}

	return nil
}

// getManifestsWithRef returns manifest list along with image reference
// Manifets are sorted by creation time descending.
func getManifestsWithRef(refsToManifest map[string]google.ManifestInfo) []Manifest {
	manifests := []Manifest{}
	for ref, manifest := range refsToManifest {
		manifests = append(manifests, Manifest{
			ref:          ref,
			manifestInfo: manifest,
		})
	}

	return byCreationTimeDesc(manifests)
}

// byCreationTimeDesc sorts manifests by creation time descending
func byCreationTimeDesc(manifests []Manifest) []Manifest {
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].manifestInfo.Created.After(
			manifests[j].manifestInfo.Created,
		)
	})

	return manifests
}

// shouldDelete determines if image should be deleted
// Last numToKeep images and ones created within grace period won't be deleted.
// Tagged images won't be deleted when cleanTags is not set.
// Otherwise all tags which does not match specified patterns will be deleted.
func shouldDelete(index int, manifest Manifest, repository Repository,
	keepTags []*regexp.Regexp) bool {

	if index < repository.NumberToKeep {
		return false
	}

	graceDate := time.Now().AddDate(0, 0, -repository.GracePeriod)
	if manifest.manifestInfo.Created.After(graceDate) {
		return false
	}

	if len(manifest.manifestInfo.Tags) > 0 {
		if !repository.CleanTags {
			return false
		}

		for _, tag := range manifest.manifestInfo.Tags {
			matched := matchTag(tag, keepTags)
			if matched {
				return false
			}
		}
		return true
	}

	return true
}

// matchTag matches tag name with regexp pattern
func matchTag(tag string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		matched := pattern.MatchString(tag)
		if matched {
			return true
		}
	}
	return false
}

// getPatterns returns compiled Regexp objects based on expression strings
// Returns error when cannot compile expression
func getPatterns(stringPatterns []string) ([]*regexp.Regexp, error) {
	patterns := []*regexp.Regexp{}
	for _, pattern := range stringPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, nil
		}
		patterns = append(patterns, re)
	}
	return patterns, nil
}

// delete deletes specific image based on manifest
func deleteImage(manifest Manifest, dry bool) {
	logrus.Infof("Cleaning image with ref: %s (tags: %s)", manifest.ref,
		manifest.manifestInfo.Tags)

	if !dry {
		// TODO do delete
	}
}

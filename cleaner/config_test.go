package cleaner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	c := Config{}
	// Validate empty repositories
	assert.EqualError(t, c.Validate(), EmptyRepositories)
	// Validate empty repository name
	c.Repositories = []Repository{
		Repository{},
	}
	assert.EqualError(t, c.Validate(), EmptyRepositoryName)
}

func TestInvalidConfigParse(t *testing.T) {
	invalidYAML := []byte("invalidYaml")
	_, err := parseConfig(invalidYAML)
	assert.Error(t, err)

	invalidStructure := []byte(`
		foo:
			bar: test
	`)
	_, err = parseConfig(invalidStructure)
	assert.Error(t, err)
}

func TestValidConfigParse(t *testing.T) {
	validConfig := []byte(`
repositories:
  - name: gcr.io/test/test/
    numberToKeep: 5
    gracePeriod: 10
    cleanTags: true
    keepTags:
      - master
      - ^v(0|\d+).(0|\d+).(0|\d+)$`,
	)

	expectedConfig := Config{
		Repositories: []Repository{
			{
				Name:         "gcr.io/test/test/",
				CleanTags:    true,
				GracePeriod:  10,
				NumberToKeep: 5,
				KeepTags:     []string{"master", "^v(0|\\d+).(0|\\d+).(0|\\d+)$"},
			},
		},
	}

	actualConfig, err := parseConfig(validConfig)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, actualConfig)
}

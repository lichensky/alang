package cleaner

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config errors
const (
	EmptyRepositories   = "repositories cannot be empty"
	EmptyRepositoryName = "repository name is required"
)

// Repository config structure
type Repository struct {
	Name         string   `yaml:"name"`
	CleanTags    bool     `yaml:"cleanTags"`
	KeepTags     []string `yaml:"keepTags"`
	NumberToKeep int      `yaml:"numberToKeep"`
	GracePeriod  int      `yaml:"gracePeriod"`
}

// Config structure
type Config struct {
	Repositories []Repository `yaml:"repositories"`
}

// Validate configuration
func (c *Config) Validate() error {
	if len(c.Repositories) == 0 {
		return errors.New(EmptyRepositories)
	}
	for _, repository := range c.Repositories {
		if repository.Name == "" {
			return errors.New(EmptyRepositoryName)
		}
	}
	return nil
}

// LoadConfig from YAML file
func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	return parseConfig(content)
}

// parseConfig parses confic from bytes
func parseConfig(bytes []byte) (Config, error) {
	config := Config{}
	err := yaml.Unmarshal(bytes, &config)

	return config, err
}

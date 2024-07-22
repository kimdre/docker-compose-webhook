package config

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/creasty/defaults"

	"github.com/caarlos0/env/v11"
	"github.com/go-git/go-billy/v5"
	"github.com/go-playground/webhooks/v6/github"
	"gopkg.in/yaml.v3"
)

var ErrConfigFileNotFound = errors.New("configuration file not found in repository")

// DeployConfigMeta is the deployment configuration meta data
type DeployConfigMeta struct {
	// DeploymentConfigFilePath is the default path/regex pattern to the deployment configuration file
	// in a repository and overrides the default deployment configuration
	DeploymentConfigFilePath string `env:"DEPLOYMENT_CONFIG_FILE_NAME" envDefault:".compose-webhook.y(a)?ml"`
}

// DeployConfig is the structure of the deployment configuration file
type DeployConfig struct {
	Name                string `yaml:"name"`                                          // Name is the name of the docker-compose deployment / stack
	Reference           string `yaml:"reference" default:"refs/heads/main"`           // Reference is the Git reference to the deployment, e.g. refs/heads/main or refs/tags/v1.0.0
	DockerComposePath   string `yaml:"compose_path" default:"docker-compose.y(a)?ml"` // DockerComposePath is the path to the docker compose file
	SkipTLSVerification bool   `yaml:"skip_tls_verify" default:"false"`               // SkipTLSVerification skips the TLS verification
}

func NewDeployConfigMeta() (*DeployConfigMeta, error) {
	cfg := DeployConfigMeta{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func DefaultDeployConfig(name string) *DeployConfig {
	return &DeployConfig{
		Reference:           "/ref/heads/main",
		Name:                name,
		DockerComposePath:   "docker-compose.y(a)?ml",
		SkipTLSVerification: false,
	}
}

func (c *DeployConfig) parseConfigFile(file []byte) error {
	err := defaults.Set(c)
	if err != nil {
		return err
	}

	type Plain DeployConfig

	if err := yaml.Unmarshal(file, (*Plain)(c)); err != nil {
		return err
	}

	return nil
}

// GetDeployConfig returns either the deployment configuration from the repository or the default configuration
func GetDeployConfig(fs billy.Filesystem, event github.PushPayload) (*DeployConfig, error) {
	m, err := NewDeployConfigMeta()
	if err != nil {
		return nil, err
	}

	lastIdx := strings.LastIndex(m.DeploymentConfigFilePath, "/")

	var path, file string

	if lastIdx == -1 {
		path = ""
		file = m.DeploymentConfigFilePath
	} else {
		path = m.DeploymentConfigFilePath[:lastIdx]
		file = m.DeploymentConfigFilePath[lastIdx+1:]
	}

	files, err := fs.ReadDir(path)
	if err != nil {
		return DefaultDeployConfig(event.Repository.Name), err
	}

	// Search for regex pattern of DeploymentConfigFilePath in the filesystem
	for _, f := range files {
		matched, err := regexp.MatchString(file, f.Name())
		if err != nil {
			return DefaultDeployConfig(event.Repository.Name), err
		}

		if matched {
			file, err := fs.Open(path + "/" + f.Name())
			defer func(f billy.File) {
				_ = f.Close()
			}(file)

			if err != nil {
				return DefaultDeployConfig(event.Repository.Name), err
			}

			// Get contents of deploy config file
			fileContents, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			c := DeployConfig{}

			if err := c.parseConfigFile(fileContents); err != nil {
				return nil, err
			}

			return &c, nil
		}
	}

	return DefaultDeployConfig(event.Repository.Name), fmt.Errorf("%s: '%s', using default configuration", m.DeploymentConfigFilePath, ErrConfigFileNotFound)
}

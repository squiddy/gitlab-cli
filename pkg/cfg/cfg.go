package cfg

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type ConfigGitlab struct {
	Group   int
	Token   string
	BaseURL string
}

type Config struct {
	GitLab ConfigGitlab
}

func getConfigFileName() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "gitlab-cli.toml"), nil
}

func LoadConfig(cfg *Config) error {
	filename, err := getConfigFileName()
	if err != nil {
		return err
	}

	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (c Config) Save() error {
	filename, err := getConfigFileName()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(c); err != nil {
		return err
	}

	return nil
}

func (c Config) Dump() {
	if err := toml.NewEncoder(os.Stdout).Encode(c); err != nil {
		return
	}
}

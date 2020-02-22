package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/squiddy/gitlab-cli/pkg/cfg"
	"github.com/xanzy/go-gitlab"
)

var config cfg.Config
var gitlabClient *gitlab.Client

func RequireConfig(cmd *cobra.Command, args []string) error {
	if config.GitLab.Token == "" {
		return errors.New("gitlab token missing")
	}
	if config.GitLab.Group == 0 {
		return errors.New("gitlab group id missing")
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use: "gitlab-cli",
}

func Execute() error {
	if err := cfg.LoadConfig(&config); err != nil {
		return err
	}

	gitlabClient = gitlab.NewClient(nil, config.GitLab.Token)
	gitlabClient.SetBaseURL(config.GitLab.BaseURL)

	return rootCmd.Execute()
}

package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup cli",
	RunE:  configure,
}

func configure(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		config.Dump()
	} else if len(args) == 2 {
		switch args[0] {
		case "gitlab.url":
			config.GitLab.BaseURL = args[1]
		case "gitlab.token":
			config.GitLab.Token = args[1]
		case "gitlab.group":
			v, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			config.GitLab.Group = v
		}
		config.Save()
	}
	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)
}

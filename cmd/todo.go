package cmd

import (
	"fmt"
	"sort"

	"github.com/dustin/go-humanize"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var todoCmd = &cobra.Command{
	Use:     "todo",
	Short:   "List of pending todos",
	PreRunE: RequireConfig,
	Args:    cobra.NoArgs,
	RunE:    ListTodos,
}

func ListTodos(cmd *cobra.Command, args []string) error {
	opts := &gitlab.ListTodosOptions{
		State: gitlab.String("pending"),
	}
	todos, _, err := gitlabClient.Todos.ListTodos(opts)
	if err != nil {
		return err
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.Before(*todos[j].CreatedAt)
	})

	fmt.Println()
	for _, t := range todos {
		fmt.Print(ansi.Green)
		fmt.Printf("    %-20v %v", humanize.Time(*t.CreatedAt), ansi.Reset)
		fmt.Printf("%-40v", t.Author.Name)
		fmt.Printf("%v", t.Target.Title)
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(todoCmd)
}

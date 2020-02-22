package cmd

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/squiddy/gitlab-cli/pkg/review"
	"github.com/xanzy/go-gitlab"
)

var reviewCmd = &cobra.Command{
	Use:     "review",
	Short:   "Working with open reviews",
	PreRunE: RequireConfig,
	Args:    cobra.NoArgs,
	RunE:    ListReviews,
}

var reviewOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open review in browser",
	RunE:  OpenReview,
}

func header(title string) {
	fmt.Printf("\n%s\n\n", ansi.Color(title, "white+b"))
}

func printReviews(reviews []review.Review) {
	for _, r := range reviews {
		fmt.Print(ansi.Green)
		fmt.Printf("    #%-3v %v", r.MergeRequest.IID, ansi.Reset)
		fmt.Printf("%-50v", r.MergeRequest.Title)
		fmt.Printf("%v[%s/%s]%v", ansi.LightBlue, r.MergeRequest.Author.Username, r.MergeRequest.SourceBranch, ansi.Reset)
		fmt.Println()

		fmt.Print("         - ")
		switch r.CheckStatus() {
		case "pending":
			fmt.Print(ansi.LightBlack, "Checks pending", ansi.Reset)
		case "unknown":
			fmt.Print(ansi.LightBlack, "Checks unknown", ansi.Reset)
		case "running":
			fmt.Print(ansi.Cyan, "Checks running", ansi.Reset)
		case "success":
			fmt.Print(ansi.Green, "Checks passing", ansi.Reset)
		case "failed":
			fmt.Print(ansi.Red, "Checks failed", ansi.Reset)
		}

		fmt.Print(" - ")
		switch r.Status() {
		case review.ReviewNeeded:
			fmt.Print(ansi.Yellow, "Review needed", ansi.Reset)
		case review.ReviewShipIt:
			fmt.Print(ansi.Green, "Review passed", " (+", r.MergeRequest.Upvotes, ")", ansi.Reset)
		}
		fmt.Println()
	}
}

func ListReviews(cmd *cobra.Command, args []string) error {
	opts := &gitlab.ListGroupMergeRequestsOptions{
		State:   gitlab.String("opened"),
		OrderBy: gitlab.String("updated_at"),
	}
	reviews, _, err := gitlabClient.MergeRequests.ListGroupMergeRequests(config.GitLab.Group, opts)
	if err != nil {
		return err
	}

	currentUser, _, err := gitlabClient.Users.CurrentUser()
	if err != nil {
		return err
	}

	var myReviews []review.Review
	var otherReviews []review.Review

	for _, r := range reviews {
		review := review.New(gitlabClient, r)
		if r.Author.ID == currentUser.ID {
			myReviews = append(myReviews, review)
		} else {
			otherReviews = append(otherReviews, review)
		}
	}

	header("Your reviews")
	printReviews(myReviews)
	header("Other reviews")
	printReviews(otherReviews)

	return nil
}

func OpenReview(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	reviews, _, err := gitlabClient.MergeRequests.ListGroupMergeRequests(config.GitLab.Group, nil)
	if err != nil {
		return err
	}

	for _, r := range reviews {
		if r.IID == id {
			if err := exec.Command("xdg-open", r.WebURL).Start(); err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	reviewCmd.AddCommand(reviewOpenCmd)
	rootCmd.AddCommand(reviewCmd)
}

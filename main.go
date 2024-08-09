package main

import (
	"os"

	"github.com/mateothegreat/go-gister/commands"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "gister",
	Short: "GitHub Gist CLI tool.",
	Long:  "Gister is a CLI tool for creating GitHub Gists.",
}

func init() {
	root.PersistentFlags().BoolP("dry-run", "", false, "Dry run the command.")
	root.PersistentFlags().StringP("username", "u", "", "GitHub username.")
	root.PersistentFlags().StringP("token", "t", os.Getenv("GITHUB_TOKEN"), "GitHub token.")
}

func main() {
	root.AddCommand(commands.Create)
	root.Execute()
}

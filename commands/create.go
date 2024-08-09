package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/v63/github"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/spf13/cobra"
)

type File struct {
	Path      string
	Filename  string
	Directory string
	Content   string
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "Create a new GitHub Gist.",
	Long:  "Create a new GitHub Gist.",
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			multilog.Fatal("create", "failed to get username", map[string]interface{}{
				"error": err,
			})
		}

		token, err := cmd.Flags().GetString("token")
		if err != nil {
			multilog.Fatal("create", "failed to get token", map[string]interface{}{
				"error": err,
			})
		}

		public, err := cmd.Flags().GetBool("public")
		if err != nil {
			multilog.Fatal("create", "failed to get public", map[string]interface{}{
				"error": err,
			})
		}

		description, err := cmd.Flags().GetString("description")
		if err != nil {
			multilog.Fatal("create", "failed to get description", map[string]interface{}{
				"error": err,
			})
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			multilog.Fatal("create", "failed to get path", map[string]interface{}{
				"error": err,
			})
		}

		files := make([]File, 0)
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				files = append(files, File{
					Path:     path,
					Filename: filepath.Base(path),
					Content:  string(content),
					Directory: func() string {
						dir := filepath.Dir(path)
						if dir == "./" {
							return ""
						}
						return dir
					}(),
				})
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error reading files:", err)
			return
		}

		gistFiles := make(map[github.GistFilename]github.GistFile)
		for _, file := range files {
			gistFiles[github.GistFilename(file.Filename)] = github.GistFile{
				Filename: github.String(file.Filename),
				Content:  github.String(file.Content),
			}
			multilog.Debug("create", "file added", map[string]interface{}{
				"filename": file.Filename,
				"path":     file.Path,
				"size":     len(file.Content),
			})
		}

		client := github.NewClient(nil).WithAuthToken(token)
		gist, _, err := client.Gists.Create(context.Background(), &github.Gist{
			Public:      github.Bool(public),
			Description: github.String(description),
			Owner:       &github.User{Login: github.String(username)},
			Files:       gistFiles,
		})
		if err != nil {
			log.Fatal(err)
		}
		multilog.Info("create", "gist created", map[string]interface{}{
			"url":   *gist.HTMLURL,
			"files": len(gistFiles),
		})
	},
}

func init() {
	multilog.RegisterLogger(multilog.LogMethod("console"), multilog.NewConsoleLogger(&multilog.NewConsoleLoggerArgs{
		Level:  multilog.DEBUG,
		Format: multilog.FormatText,
	}))

	Create.Flags().StringP("username", "u", "", "GitHub username.")
	Create.MarkFlagRequired("username")

	Create.Flags().StringP("token", "t", "", "GitHub token.")
	Create.MarkFlagRequired("token")

	Create.Flags().StringP("description", "d", "", "Gist description.")
	Create.MarkFlagRequired("description")

	Create.Flags().StringP("path", "", "", "Path to the file or directory to gist.")
	Create.MarkFlagRequired("path")

	Create.Flags().BoolP("public", "p", false, "Make the gist public.")
}

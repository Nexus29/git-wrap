package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Nexus29/git-wrap/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupQuestions = []*survey.Question{
	{Name: "username", Prompt: &survey.Input{Message: "Enter your GitHub username:"}},
	{Name: "token", Prompt: &survey.Password{Message: "Enter your GitHub Personal Access Token (PAT):"}},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure your GitHub credentials for git-wrap",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚙️  Setting up git-wrap configuration profile...")

		answers := struct {
			Username string
			Token    string
		}{}

		if err := survey.Ask(setupQuestions, &answers); err != nil {
			fmt.Println("❌ Setup cancelled.")
			return
		}

		cfg := &config.Config{
			GitHubUsername: strings.TrimSpace(answers.Username),
			GitHubToken:    strings.TrimSpace(answers.Token),
		}

		if err := config.Save(cfg); err != nil {
			fmt.Println("❌ Failed to save configuration file:", err)
			return
		}

		fmt.Println("✨ Setup complete! Credentials securely stored in ~/.git-wrap.json. You can now use 'git-wrap save'.")
	},
}

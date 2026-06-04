package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(submoduleCmd)
}

var submoduleCmd = &cobra.Command{
	Use:   "submodule",
	Short: "Interactive utility to inject or attach submodules cleanly",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧩 --- Smart Submodule Injection Wizard ---")

		repoURL := ""
		urlPrompt := &survey.Input{
			Message: "Enter the remote Git submodule repository URL:",
		}
		if err := survey.AskOne(urlPrompt, &repoURL); err != nil || strings.TrimSpace(repoURL) == "" {
			fmt.Println("❌ Submodule target address required.")
			return
		}

		targetPath := ""
		pathPrompt := &survey.Input{
			Message: "Enter the local target path directory (e.g., libs/my-sub):",
		}
		if err := survey.AskOne(pathPrompt, &targetPath); err != nil || strings.TrimSpace(targetPath) == "" {
			fmt.Println("❌ Target installation path required.")
			return
		}

		fmt.Printf("📦 Attaching submodule link: %s -> into path: %s\n", repoURL, targetPath)
		addCmd := exec.Command("git", "submodule", "add", strings.TrimSpace(repoURL), strings.TrimSpace(targetPath))
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr

		if err := addCmd.Run(); err != nil {
			fmt.Println("❌ Failed to add submodule cleanly.")
			return
		}

		fmt.Println("🚀 Submodule attached, initialized, and integrated successfully!")
	},
}

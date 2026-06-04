package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Nexus29/git-wrap/pkg/config"
	"github.com/Nexus29/git-wrap/pkg/submodules"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(saveCmd)
}

var commitQuestions = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Select the type of change you are committing:",
			Options: []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore"},
		},
	},
	{
		Name: "scope",
		Prompt: &survey.Input{
			Message: "Enter the scope of this change (optional):",
		},
	},
	{
		Name: "subject",
		Prompt: &survey.Input{
			Message: "Write a short, imperative description (lowercase, no period):",
		},
	},
}

type CreateRepoPayload struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Validates repo state, auto-creates on GitHub if missing, stages, tracks submodules, and pushes code",
	Run: func(cmd *cobra.Command, args []string) {
		
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("❌ Configuration not found! Please run 'git-wrap setup' first to configure your credentials.")
			return
		}

		// 1. Check for local repository structure tracking
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			fmt.Println("⚠️  No Git repository found in this directory!")

			repoName := ""
			prompt := &survey.Input{Message: "Enter your remote repository name:"}
			if err := survey.AskOne(prompt, &repoName); err != nil || strings.TrimSpace(repoName) == "" {
				fmt.Println("❌ A valid repository name is required.")
				return
			}
			repoName = strings.TrimSpace(repoName)

			isPrivate := false
			privatePrompt := &survey.Confirm{Message: "Do you want this repository to be Private?", Default: true}
			survey.AskOne(privatePrompt, &isPrivate)

			// --- AUTOMATED GITHUB API REPOSITORY PROVISIONING ---
			fmt.Printf("🌐 Instructing GitHub API to create '%s/%s'...\n", cfg.GitHubUsername, repoName)
			payload := CreateRepoPayload{Name: repoName, Private: isPrivate}
			jsonPayload, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "https://api.github.com/user/repos", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Authorization", "Bearer "+cfg.GitHubToken)
			req.Header.Set("Accept", "application/vnd.github+json")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil || (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK) {
				fmt.Printf("❌ GitHub API rejected creation or repository already exists.\n")
			} else {
				fmt.Println("🚀 GitHub repository successfully generated remotely!")
			}
			if resp != nil {
				resp.Body.Close()
			}

			fmt.Println("🔄 Initializing a fresh local Git repository...")
			exec.Command("git", "init").Run()
			exec.Command("git", "branch", "-M", "main").Run()

			remoteURL := fmt.Sprintf("git@github.com:%s/%s.git", cfg.GitHubUsername, repoName)
			fmt.Printf("🔗 Setting remote origin to: %s\n", remoteURL)
			exec.Command("git", "remote", "add", "origin", remoteURL).Run()
		}

		// 2. Call the new Submodule Manipulation package engine
		submodules.ProcessSubmodules()

		// 3. Standard Staging Routine Flow
		fmt.Println("🟢 Staging files via 'git add .' ...")
		exec.Command("git", "add", ".").Run()

		// 4. Launch European Commission UI guidelines wizard
		fmt.Println("\n--- 🇪🇺 European Commission Commit Wizard ---")
		answers := struct {
			Type    string
			Scope   string
			Subject string
		}{}

		if err := survey.Ask(commitQuestions, &answers); err != nil {
			fmt.Println("❌ Interrupted wizard generation scope context:", err)
			return
		}

		subject := strings.ToLower(strings.TrimSpace(answers.Subject))
		if strings.HasSuffix(subject, ".") {
			subject = strings.TrimSuffix(subject, ".")
		}

		var commitMsg string
		scope := strings.TrimSpace(answers.Scope)
		if scope != "" {
			commitMsg = fmt.Sprintf("%s(%s): %s", answers.Type, scope, subject)
		} else {
			commitMsg = fmt.Sprintf("%s: %s", answers.Type, subject)
		}

		fmt.Printf("\nGenerated EC Commit Message:\n> %s\n\n", commitMsg)

		// 5. Local Committing Routines Execution
		fmt.Println("💾 Committing changes...")
		commitExec := exec.Command("git", "commit", "-m", commitMsg)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr
		if err := commitExec.Run(); err != nil {
			fmt.Println("❌ Commit skipped or no current tracking alterations verified.")
			return
		}

		// 6. Track Active Working Branch Identity
		branchBytes, err := exec.Command("git", "branch", "--show-current").Output()
		currentBranch := "main"
		if err == nil && len(branchBytes) > 0 {
			currentBranch = strings.TrimSpace(string(branchBytes))
		}

		// 7. Final Push Payload Assembly Execution 
		fmt.Printf("🚀 Pushing code to GitHub (origin %s)...\n", currentBranch)
		pushExec := exec.Command("git", "push", "-u", "origin", currentBranch)
		pushExec.Stdout = os.Stdout
		pushExec.Stderr = os.Stderr
		if err := pushExec.Run(); err != nil {
			fmt.Println("❌ Failed pushing upstream elements safely.")
			return
		}

		fmt.Printf("\n🎉 Workflow successfully completed! Your code is live on GitHub under %s/%s\n", cfg.GitHubUsername, currentBranch)
	},
}

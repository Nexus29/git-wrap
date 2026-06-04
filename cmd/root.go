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
	"github.com/spf13/cobra"
)

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

var rootCmd = &cobra.Command{
	Use:   "git-wrap",
	Short: "git-wrap is an EC-compliant Git and Submodule manager",
	// Adding the Run block here ensures the tool runs immediately instead of showing a help menu!
	Run: func(cmd *cobra.Command, args []string) {
		
		// 1. Check for Credentials
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("⚙️  Welcome to git-wrap! Setting up your configuration profile...")
			setupCredentials()
			return
		}

		// 2. Check for Local Repo Context Tracking
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

		// 3. --- SUBMODULE MANIPULATION ENGINE ---
		processSubmodules()

		// 4. Standard Staging Routine Flow
		fmt.Println("🟢 Staging files via 'git add .' ...")
		exec.Command("git", "add", ".").Run()

		// 5. Launch European Commission UI guidelines wizard
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

		// 6. Local Committing Routines Execution
		fmt.Println("💾 Committing changes...")
		commitExec := exec.Command("git", "commit", "-m", commitMsg)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr
		if err := commitExec.Run(); err != nil {
			fmt.Println("❌ Commit skipped or no current tracking alterations verified.")
			return
		}

		// 7. Track Active Working Branch Identity
		branchBytes, err := exec.Command("git", "branch", "--show-current").Output()
		currentBranch := "main"
		if err == nil && len(branchBytes) > 0 {
			currentBranch = strings.TrimSpace(string(branchBytes))
		}

		// 8. Final Push Payload Execution 
		fmt.Printf("🚀 Pushing code to GitHub (origin %s)...\n", currentBranch)
		pushExec := exec.Command("git", "push", "-u", "origin", currentBranch)
		pushExec.Stdout = os.Stdout
		pushExec.Stderr = os.Stderr
		if err := pushExec.Run(); err != nil {
			fmt.Println("❌ Failed pushing upstream elements safely.")
			return
		}

		fmt.Printf("\n🎉 Workflow successfully completed! Your code is live under %s/%s\n", cfg.GitHubUsername, currentBranch)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupCredentials() {
	var setupQuestions = []*survey.Question{
		{Name: "username", Prompt: &survey.Input{Message: "Enter your GitHub username:"}},
		{Name: "token", Prompt: &survey.Password{Message: "Enter your GitHub Personal Access Token (PAT):"}},
	}
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
	config.Save(cfg)
	fmt.Println("✨ Setup complete! Run the script again to execute the loop.")
}

func processSubmodules() {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return // No submodules to handle
	}

	fmt.Println("🔍 Inspecting submodule architecture pipelines...")
	
	exec.Command("git", "submodule", "sync").Run()
	exec.Command("git", "submodule", "update", "--init", "--recursive").Run()

	out, err := exec.Command("git", "config", "--file", ".gitmodules", "--get-regexp", "path").Output()
	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		submodulePath := parts[1]

		fmt.Printf("📂 Checking logs for submodule path: %s\n", submodulePath)

		fetchCmd := exec.Command("git", "fetch")
		fetchCmd.Dir = submodulePath
		fetchCmd.Run()

		// Checks if upstream branch has new updates compared to our current pointer
		logCmd := exec.Command("git", "log", "..origin/HEAD", "--oneline")
		logCmd.Dir = submodulePath
		logBytes, err := logCmd.Output()

		if err == nil && len(logBytes) > 0 {
			logStr := string(logBytes)
			// Target the '[track-update]' keyword inside incoming submodule commits
			if strings.Contains(logStr, "[track-update]") {
				fmt.Printf("🔥 Found '[track-update]' tag inside submodule '%s'! Fast-forwarding...\n", submodulePath)
				
				mergeCmd := exec.Command("git", "merge", "origin/HEAD", "--ff-only")
				mergeCmd.Dir = submodulePath
				if err := mergeCmd.Run(); err != nil {
					fmt.Printf("⚠️  Could not automatically fast-forward submodule '%s'.\n", submodulePath)
				}
			}
		}
	}
}

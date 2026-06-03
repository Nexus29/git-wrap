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
	rootCmd.AddCommand(saveCmd)
}

// Interactive prompt schema for the EC Commit Message Guidelines
var commitQuestions = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Select the type of change you are committing:",
			Options: []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore"},
			Help:    "feat: New feature | fix: Bug resolution | docs: Documentation adjustments",
		},
	},
	{
		Name: "scope",
		Prompt: &survey.Input{
			Message: "Enter the scope of this change (optional, e.g., ui, api, core):",
		},
	},
	{
		Name: "subject",
		Prompt: &survey.Input{
			Message: "Write a short, imperative description (lowercase, no trailing period):",
		},
	},
}

var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Validates repo state, stages, runs EC wizard, and pushes code",
	Run: func(cmd *cobra.Command, args []string) {
		
		// 1. Check for Git Repository tracking
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			fmt.Println("⚠️  No Git repository found in this directory!")
			
			// Initialize local repo
			fmt.Println("🔄 Initializing a fresh local Git repository...")
			if err := exec.Command("git", "init").Run(); err != nil {
				fmt.Println("❌ Failed to initialize local git repo:", err)
				return
			}

			// Rename default branch to main
			exec.Command("git", "branch", "-M", "main").Run()

			// Prompt for remote repository configurations
			repoName := ""
			prompt := &survey.Input{
				Message: "Enter your remote repository name (GitHub):",
			}
			if err := survey.AskOne(prompt, &repoName); err != nil || strings.TrimSpace(repoName) == "" {
				fmt.Println("❌ A valid repository name is required.")
				return
			}

			// Hook up remote target configuration matching your GitHub profile
			remoteURL := fmt.Sprintf("git@github.com:Nexus29/%s.git", strings.TrimSpace(repoName))
			fmt.Printf("🔗 Setting remote origin to: %s\n", remoteURL)
			if err := exec.Command("git", "remote", "add", "origin", remoteURL).Run(); err != nil {
				fmt.Println("❌ Failed to bind remote origin:", err)
				return
			}
		}

		// 2. Stage current operational footprint context
		fmt.Println("🟢 Staging files via 'git add .' ...")
		if err := exec.Command("git", "add", ".").Run(); err != nil {
			fmt.Println("❌ Failed staging execution framework profiles:", err)
			return
		}

		// 3. Launch EC Guideline wizard terminal dialog inputs
		fmt.Println("\n--- 🇪🇺 European Commission Commit Wizard ---")
		answers := struct {
			Type    string
			Scope   string
			Subject string
		}{}

		if err := survey.Ask(commitQuestions, &answers); err != nil {
			fmt.Println("❌ Interrupted wizard engine context generation:", err)
			return
		}

		// 4. Formatting checks parsing standard validation components
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

		// 5. Native Git Commit Processing Execution
		fmt.Println("💾 Committing changes...")
		commitExec := exec.Command("git", "commit", "-m", commitMsg)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr
		if err := commitExec.Run(); err != nil {
			fmt.Println("❌ Commit routine optimization rejected or no changes found.")
			return
		}

		// 6. Inspect current tracking branch identity structures
		branchBytes, err := exec.Command("git", "branch", "--show-current").Output()
		currentBranch := "main"
		if err == nil && len(branchBytes) > 0 {
			currentBranch = strings.TrimSpace(string(branchBytes))
		}

		// 7. Dynamic multi-stream secure network push protocol 
		fmt.Printf("🚀 Pushing code to GitHub (origin %s)...\n", currentBranch)
		pushExec := exec.Command("git", "push", "-u", "origin", currentBranch)
		pushExec.Stdout = os.Stdout
		pushExec.Stderr = os.Stderr
		if err := pushExec.Run(); err != nil {
			fmt.Println("❌ Failed pushing upstream assets payload safely.")
			return
		}

		fmt.Println("\n🎉 Workflow successfully completed!")
	},
}

package submodules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ProcessSubmodules loops over your .gitmodules file and tracks updates
func ProcessSubmodules() {
	if _, err := os.Stat(".gitmodules"); os.IsNotExist(err) {
		return 
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

		logCmd := exec.Command("git", "log", "..origin/HEAD", "--oneline")
		logCmd.Dir = submodulePath
		logBytes, err := logCmd.Output()

		if err == nil && len(logBytes) > 0 {
			logStr := string(logBytes)
			if strings.Contains(logStr, "(track-update)") {
				fmt.Printf("🔥 Found 'track-update' tag inside submodule '%s'! Fast-forwarding...\n", submodulePath)
				
				mergeCmd := exec.Command("git", "merge", "origin/HEAD", "--ff-only")
				mergeCmd.Dir = submodulePath
				if err := mergeCmd.Run(); err != nil {
					fmt.Printf("⚠️  Could not automatically fast-forward submodule '%s'.\n", submodulePath)
				}
			}
		}
	}
}

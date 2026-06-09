# git-wrap 🚀

`git-wrap` is a lightning-fast Git CLI wrapper written in Go. It is designed to optimize daily development workflows by combining local repository checks, multi-user configuration onboarding, automated submodule lifecycle tracking, and an interactive commit wizard enforcing the strict **European Commission (EC) Git Commit Conventions**.

---

## ✨ Features

* **⚙️ One-Time Shared Setup (`git-wrap setup`)**: Prompts for and securely stores a user's GitHub username and Personal Access Token (PAT) inside `~/.git-wrap.json`. This makes the tool fully portable and shareable with other developers!
* **🌱 Zero-Config Auto-Initialization**: If running inside an untracked directory, `git-wrap` automatically runs `git init`, prompts for the remote repository name and privacy status, provisions the repo directly on GitHub via its API, and binds the local workspace to the new origin.
* **🔍 Automated Submodule Synchronization**: Intelligently scans incoming logs and commit scopes for registered submodules. If the structural keyword track (`(track-update)`) is parsed, it automatically fast-forwards the pointer and registers it to the parent state.
* **🟢 Unified Workflow (`git-wrap save`)**: Chains repository evaluation, submodule tracking updates, `git add .`, interactive EC-compliant commit generation, and `git push` into a single, seamless command execution loop.
* **🛡️ Built in Go with Cobra**: Completely modular architecture powered by the Cobra CLI framework. Compiles down to a single, cross-platform binary with zero external runtime environments required.

---

## 📂 Project Architecture

The architecture maintains strict modular decoupling, splitting command routers from backend utility logics:

```text
git-wrap/
├── cmd/
│   ├── root.go          # Core Cobra CLI base framework routing
│   ├── save.go          # Core 'save' workflow engine (Repo init, API calls & Wizard)
│   ├── setup.go         # Command router to capture user GitHub tokens
│   └── submodule.go     # Interactive wizard to inject fresh submodules
├── pkg/
│   ├── config/
│   │   └── config.go    # Safe management & serialization of ~/.git-wrap.json
│   └── submodules/
│       └── track.go     # Isolated logic for log fetching & pointer manipulation
├── go.mod               # Go module dependencies declaration
├── go.sum               # Cryptographic checksums for packages
├── PKGBUILD             # Arch Linux native package deployment blueprint
└── main.go              # Simple application execution entrypoint
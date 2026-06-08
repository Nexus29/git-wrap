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

```

---

## 🛠️ Requirements & Installation

### Prerequisites

* **Git**: Available globally in your system path environment.
* **Go**: Version `1.21` or higher (to compile from source).

### 📦 Installing on Arch Linux (Pacman)

The project includes native support for Arch Linux using `makepkg`. It is configured to compile your binary securely in isolation using your global temporary system directories without dirtying or overwriting your project's Go source `pkg/` files.

```bash
# 1. Compile the snapshot binary safely
goreleaser release --snapshot --clean

# 2. Build and install natively via pacman
makepkg -si

```

To cleanly uninstall the utility from your Arch Linux system:

```bash
sudo pacman -R git-wrap

```

### 🛠️ Installing From Source (Alternative)

```bash
# Clone the repository
git clone [https://github.com/Nexus29/git-wrap.git](https://github.com/Nexus29/git-wrap.git)
cd git-wrap

# Fetch dependencies and verify modules
go mod download
go mod verify

# Compile and install globally in your $GOPATH/bin
go install

```

---

## 🚀 Quick Start & Usage

### 1. The First-Time Setup

Before using the automation engine, users must configure their local environment so the tool knows where to look and push. Run:

```bash
git-wrap setup

```

This will securely prompt you for your **GitHub Username** and a **Personal Access Token (PAT)** (Classic PAT with `repo` scope enabled), storing them locally in `~/.git-wrap.json`.

### 2. Execute the Unified Save Sequence

Run the core lifecycle tool in any active working directory:

```bash
git-wrap save

```

This triggers the sequential pipeline:

1. **Onboarding / Repo Evaluation**: Checks for a local `.git` structure. If missing, it creates it, leverages your configured credentials to automatically make a public/private repo on GitHub via the API, and links your local origin to `git@github.com:<your-username>/<repo-name>.git`.
2. **Submodule Check**: Loops over your `.gitmodules` file, triggers a remote `git fetch` inside every submodule, and inspects the remote logs and commit metadata for tracking updates.
3. **Staging**: Automatically triggers `git add .` to capture file updates and any fast-forwarded submodule links.
4. **EC Commit Message Wizard**: Launches the interactive prompt flow.
5. **Push**: Executes `git commit` and pushes code live to your upstream branch.

### 3. Interactive Submodule Management

To add a brand new submodule cleanly to your workspace layout, execute:

```bash
git-wrap submodule

```
---

## 🔄 How Automated Submodule Tracking Works

The submodule manipulation tracker is fully automated and relies on explicit scope tags within your code history.

1. **Tagging updates**: When committing changes meant to sync, inclusion of the **`(track-update)`** tag inside the commit scope registers with the parser.
2. **Automated Fetch & Sync**: When `git-wrap save` runs, the parsing engine scans the target scopes of incoming logs. If a matching update track is detected, it executes a fast-forward merge (`git merge origin/HEAD --ff-only`) directly inside the target submodule workspace directory, cleanly bumping the pointer before final parent staging.

---

## ⚙️ European Commission Commit Conventions Reference

The framework strictly structures layout standards following the [European Commission Git Guidelines](https://ec.europa.eu/component-library/v1.15.0/eu/docs/conventions/git/):

* **Types**: Must align precisely with target definitions specifying modification intents (`feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`).
* **Scope**: Wrapping brackets denoting contextual domains impacted by the code footprint (e.g., `(core)`, `(ui)`, `(api)`, `(track-update)`).
* **Subject**: Explicitly imperative present-tense framing describing the action footprint. Must begin with lowercase parameters and avoid closing punctuality attributes (`.`).

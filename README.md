# git-wrap 🚀

`git-wrap` is a lightning-fast, zero-dependency Git CLI wrapper written in Go. It is designed to optimize daily development workflows by combining local repository checks, staging, automated submodule lifecycle tracking, and an interactive commit wizard enforcing the strict **European Commission (EC) Git Commit Conventions** (based on the Conventional Commits specification).

This tool serves as an excellent foundational layer for developers managing multi-repository architectures, setting up automated environments, or laying down system-level tools for advanced, automated Linux installations.

---

## ✨ Features

* **🌱 Zero-Config Auto-Initialization**: If running inside a directory that is not yet tracked by Git (no `.git` directory found), `git-wrap` automatically runs `git init`, prompts you for a remote repository name, handles the remote origin provisioning for your account, and tracks the main upstream branch smoothly.
* **📦 Unified Workflow (`git-wrap save`)**: Chain repository validation, `git add .`, interactive commit generation, and `git push` into a single, seamless command execution loop.
* **🇪🇺 EC Commit Compliance Wizard**: Built-in interactive CLI prompts ensuring your commit structures fully respect the European Commission component library guidelines: `<type>(<scope>): <subject>`.
* **🧩 Smart Submodule Embedding**: On-the-fly interactive wizard options to cleanly inject, init, and recursively fetch remote Git submodules.
* **🔍 Automated Keyword Track Integration**: Intelligently scans incoming submodule logs during execution. If a specific structural keyword track (e.g., `[track-update]`) is parsed from upstream submodules, it automatically fast-forwards your pointers and bumps your parent state seamlessly.
* **🛡️ Built in Go**: Compiles to a single, cross-platform binary. Zero external runtime environments (NodeJS, Python, Ruby) are required on the host system.

---

## 📂 Project Architecture

The architecture maintains strict modular decoupling, splitting command routers from low-level systems commands wrappers:

```text
git-wrap/
├── cmd/
│   ├── root.go          # Core CLI configuration routing
│   ├── save.go          # The primary 'save' workflow engine, Repository Init, & EC Wizard
│   └── submodule.go     # Interactive submodule attachment controllers
├── pkg/
│   ├── git/
│   │   └── git.go       # System wrappers executing native Git interactions
│   └── submodules/
│       └── track.go     # Parsing logs & keyword evaluations
├── go.mod               # Go module dependencies declaration
├── go.sum               # Cryptographic checksums for packages
└── main.go              # Execution application entrypoint

```

---

## 🛠️ Requirements & Installation

### Prerequisites

* **Git**: Your system must have standard `git` binary distributions available globally in your path environment.
* **Go** (Only required to compile from source): version `1.18` or higher.

### Installing From Source

To install `git-wrap` directly onto your local Linux, macOS, or Windows subsystem, follow these steps:

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

*Ensure your environment path configuration contains `$GOPATH/bin` (typically `~/go/bin`) to invoke the command globally.*

---

## 🚀 Quick Start & Usage

### 1. Execute the Unified Save Sequence

Run the core lifecycle tool in any active directory:

```bash
git-wrap save

```

This triggers a sequential pipeline:

1. **Repository Evaluation**: Checks for a local `.git` structure. If missing, it safely initializes the folder, asks for your target remote repository name, links it to your GitHub account (`git@github.com:Nexus29/<repo-name>.git`), and prepares your initial branch context.
2. Automatically triggers `git add .` to stage your ongoing context files.
3. Launches the interactive **EC Commit Message Wizard**.
4. Scans tracking rules for external submodules.
5. Performs isolated `git commit -m "<structured-message>"` generation.
6. Injects changes into your remote origin branch dynamically.

### 2. The Interactive Commit Protocol

When executing the wizard, `git-wrap` guides your framing to enforce valid structural definitions:

```text
? Select the type of change you are committing: (Use arrow keys)
❯ feat      (Introducing new features)
  fix       (Bug resolutions)
  docs      (Documentation changes only)
  style     (Formatting changes, missing semi-colons, etc.)
  refactor  (Code changes that neither fix a bug nor add a feature)
  perf      (Performance improvements)
  test      (Adding missing tests or correcting existing tests)
  chore     (Updating build tasks, package configurations, etc.)

? Enter the scope of this change (optional, e.g., ui, api): core
? Write a short, imperative tense description (lowercase, no period): integrate submodule tracker

```

**Resulting Standard Output Message**: `feat(core): integrate submodule tracker`

---

## ⚙️ European Commission Commit Conventions Reference

The framework strictly structures layout standards following the [European Commission Git Guidelines](https://ec.europa.eu/component-library/v1.15.0/eu/docs/conventions/git/):

* **Types**: Must align precisely with target definitions specifying modification intents (`feat`, `fix`, etc.).
* **Scope**: Wrapping brackets denoting contextual domains impacted by the code footprint (`(ui)`, `(api)`).
* **Subject**: Explicitly imperative present-tense framing describing the action footprint. Must begin with lowercase parameters and avoid closing punctuality attributes (`.`).

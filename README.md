# Mob

<p align="center">
  <img src="https://media.giphy.com/media/3o7TKSjRrfIPjeiVyM/giphy.gif" alt="Mop cleaning up your messy code" width="300"/>
  <br>
  <em>Cleaning up your messy commits</em>
</p>

A CLI tool that helps developers manage their Git workflow and produce high-quality Pull Requests. Mob streamlines the process of creating feature branches from GitHub issues, tracking commits, and reviewing code with AI-powered suggestions based on LLVM Coding Standards.

## Features

- **Branch Management** - Create structured `wip/<issue>` and `pr/<issue>` branches linked to GitHub issues
- **Commit Tracking** - Automatically track fork points and squash commits when updating PRs
- **Interactive Review UI** - Terminal UI with syntax-highlighted diffs and customizable checklist
- **AI Code Review** - Get code suggestions powered by OpenAI GPT models

## Installation

### Dependencies

- [Go](https://golang.org/) 1.25 or higher
- [Git](https://git-scm.com/)
- [GitHub CLI](https://cli.github.com/) (`gh`) - for fetching GitHub issues

### Build from Source

```bash
git clone https://github.com/joaosaffran/mob.git
cd mob

# Build for your platform
make build-windows  # Windows
make build-linux    # Linux
make build-mac      # macOS Intel
make build-mac-arm  # macOS Apple Silicon
```

The binary will be created in the `bin/` directory.

## Commands

### init

Creates a new work-in-progress branch from a GitHub issue.

```bash
mob init
```

This displays an interactive list of GitHub issues assigned to you. Select an issue to create a `wip/<issue-number>` branch.

**Options:**

```bash
mob init --base-branch main    # Specify the base branch
mob init -b develop            # Short form
```

### update

Squashes commits from your `wip/<issue>` branch into a `pr/<issue>` branch and pushes to remote.

```bash
mob update -m "Add user authentication feature"
```

This creates or updates the `pr/<issue>` branch with a single squashed commit containing all changes since the fork point.

### review

Opens an interactive terminal UI to review your changes before creating a PR.

```bash
mob review
```

The review UI displays:

- **Left Panel** - Syntax-highlighted diff of your changes
- **Top Right** - Checklist items from `.mob/checklist.yaml`
- **Bottom Right** - AI recommendations (requires `OPENAI_API_KEY` environment variable)

**Navigation:**

| Key | Action |
|-----|--------|
| `Tab` | Switch between panels |
| `↑/↓` | Navigate items |
| `Space/Enter` | Toggle checklist / View recommendation |
| `q` | Quit |

**Configuration:**

Create `.mob/checklist.yaml` in your repository:

```yaml
items:
  - description: "Code follows project style guidelines"
  - description: "Tests are included and passing"
  - description: "Documentation is updated"
```

Set your OpenAI API key for AI recommendations:

```bash
# Windows (cmd)
set OPENAI_API_KEY=your-api-key

# Windows (PowerShell)
$env:OPENAI_API_KEY="your-api-key"

# Linux/macOS
export OPENAI_API_KEY="your-api-key"
```

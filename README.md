<a id="readme-top"></a>

# sh_adow

<!-- PROJECT LOGO -->
<a href="https://github.com/chhlga/sh_adow">
  <img src="assets/blobashik_shadow.png" alt="Bashik Shadow Mascot" width="254" height="254">
</a>


<div>
  <p>
    Simple file versioning without git complexity
    <br />
    <a href="#getting-started"><strong>Get Started</strong></a>
    ·
    <a href="#usage">View Examples</a>
    ·
    <a href="https://github.com/chhlga/sh_adow/issues">Report Bug</a>
    ·
    <a href="https://github.com/chhlga/sh_adow/issues">Request Feature</a>
  </p>
</div>

<!-- PROJECT SHIELDS -->
  [![Go Version][go-shield]][go-url]
  [![License][license-shield]][license-url]
  ![Cobra][cobra-badge]
  ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/chhlga/sh_adow/go.yml?style=for-the-badge)

---

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#what-im-looking-at">What I'm Looking At?</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#install">Install</a></li>
    <li><a href="#run">Run</a></li>
    <li><a href="#built-with">Built With</a></li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#configuration">Configuration</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

---

## What I'm Looking At?

**sh_adow** is a lightweight CLI tool for managing file versions without the complexity of git. Perfect for quick snapshots of config files, scripts, or any text file when you need simple versioning.

**Key benefits:**
- ✅ Zero setup - auto-initializes on first save
- ✅ Works with single files - no repository management
- ✅ Simple 4-command interface - save, list, restore, delete
- ✅ Flexible storage - local, centralized, or relative paths
- ✅ Beautiful terminal UI - built with Charm stack

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Features

- 🚀 **Auto-initialization** - No setup needed, just start saving files
- 📦 **Flat storage** - Simple `.shadow/` directory with flat file structure
- 🎯 **Flexible repo paths** - Store shadows locally, centrally, or relative
- 💅 **Beautiful TUI** - Built with Charm stack (Cobra, Huh, Lipgloss)
- 🔍 **Virtual HEAD** - Current file state always visible without extra storage
- 🏷️ **Tags & Notes** - Organize versions with multiple tags and descriptions
- ⚡ **Fast & Lightweight** - ~6.7MB binary, no runtime dependencies

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Install

```bash
# Via Go install
go install github.com/chhlga/sh_adow@latest
```

**OR** build from source:

```bash
git clone https://github.com/chhlga/sh_adow
cd sh_adow
make  # Builds to bin/shadow
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Run

```bash
# Save a version of a file
shadow save config.yaml -t "before-refactor" -n "Working config"

# List all tracked files
shadow list

# List versions of a specific file
shadow list config.yaml

# Restore a version
shadow restore config.yaml abc123

# Delete a version
shadow delete config.yaml abc123 --force
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- GETTING STARTED -->
## Getting Started

Shadow requires no initial setup or configuration. Just install the binary and start saving files.

### Prerequisites

- **Go 1.21+** (only if building from source)
- No other runtime dependencies required

### Installation

**Option 1: Go Install (Recommended)**

```bash
go install github.com/chhlga/sh_adow@latest
```

**Option 2: Build from Source**

```bash
# Clone the repository
git clone https://github.com/chhlga/sh_adow
cd sh_adow

# Build using Makefile
make

# OR build manually
go build -o bin/shadow .

# Optionally install to system
make install  # Installs to /usr/local/bin/shadow
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- USAGE EXAMPLES -->
## Usage

### Commands

#### `shadow save <file>`

Save a version of a file. Auto-creates `.shadow/` directory if needed.

```bash
# Save with tags and notes
shadow save config.yaml -t "stable" -t "v1.0" -n "Production config"

# Interactive mode (prompts for tags/notes)
shadow save config.yaml
```

#### `shadow list [file]`

List tracked files or versions of a specific file.

```bash
# Show all tracked files
shadow list

# Show versions of a specific file
shadow list config.yaml
```

#### `shadow restore <file> <version-id>`

Restore a file to a specific version.

```bash
# Interactive (prompts to save current state)
shadow restore config.yaml abc123

# Skip save prompt
shadow restore config.yaml abc123 --no-save
```

#### `shadow delete <file> <version-id>`

Delete a specific version.

```bash
# Interactive (prompts for confirmation)
shadow delete config.yaml abc123

# Force delete without prompt
shadow delete config.yaml abc123 --force
```

### Use Cases

- **Config files** - Track changes to dotfiles, app configs, etc.
- **Scripts** - Version important scripts before modifications
- **Documentation** - Save drafts while editing
- **Data files** - Quick snapshots of CSV, JSON, YAML files
- **Any text files** - No git setup required

### Example Scenarios

**Scenario 1: Editing a config file**

```bash
# Before making risky changes
shadow save ~/.config/app/config.yml -t "working" -n "Before adding new feature"

# Make changes to config.yml
vim ~/.config/app/config.yml

# Something broke? Restore easily
shadow list ~/.config/app/config.yml
shadow restore ~/.config/app/config.yml a1b2c3d4
```

**Scenario 2: Centralized backup**

Configure centralized storage:
```bash
mkdir -p ~/.config/sh_adow
echo "repo_path: \"~/.shadow_backups/\"" > ~/.config/sh_adow/config.yml
```

Now all shadows go to one place:
```bash
shadow save ~/important/file1.txt
shadow save ~/projects/app/config.yaml
shadow list  # Shows all files from ~/.shadow_backups/
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

## Configuration

Shadow works out of the box with zero configuration. Optionally, create `~/.config/sh_adow/config.yml` to customize shadow repository location:

```yaml
# Default: store .shadow/ in same directory as file
repo_path: "./"

# Centralized storage: all shadows in one location
repo_path: "~/.local/cache/"

# Relative path: shadow in parent/cache directory
repo_path: "../cache/"
```

### Configuration Examples

**Default (local storage):**
```yaml
repo_path: "./"
```
```
File: ~/proj/config.yaml
Shadow: ~/proj/.shadow/
```

**Centralized storage:**
```yaml
repo_path: "~/.local/cache/"
```
```
File: ~/proj/config.yaml
Shadow: ~/.local/cache/.shadow/

File: ~/work/settings.yaml
Shadow: ~/.local/cache/.shadow/  # Same repo!
```

**Relative path:**
```yaml
repo_path: "../cache/"
```
```
File: ~/proj/subdir/config.yaml
Shadow: ~/proj/cache/.shadow/
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Built with ❤️ using the [Charm](https://charm.sh) stack:

* [Cobra](https://github.com/spf13/cobra) - CLI framework
* [Huh](https://github.com/charmbracelet/huh) - Interactive forms
* [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
* [yaml.v3](https://github.com/go-yaml/yaml) - Config parsing

Approved by [Bashik](https://github.com/chhlga/chhlga/blob/main/bashik.md) 🐱

<p align="right">(<a href="#readme-top">back to top</a>)</p>

---

<!-- MARKDOWN LINKS & BADGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[go-shield]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white
[go-url]: https://go.dev/
[go-badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white

[license-shield]: https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge
[license-url]: LICENSE

[cobra-badge]: https://img.shields.io/badge/Cobra-CLI-blue?style=for-the-badge

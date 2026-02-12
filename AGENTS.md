# Shadow - File Versioning CLI

## Status: ✅ IMPLEMENTED

A lightweight Go CLI application for managing multiple versions of any file without git complexity.

## Implementation Summary

### What Was Built

A fully functional CLI tool with the following commands:
- `shadow save <file>` - Save file versions with tags/notes
- `shadow list [file]` - List tracked files or specific file versions
- `shadow restore <file> <version-id>` - Restore previous versions
- `shadow delete <file> <version-id>` - Delete specific versions

### Architecture Decisions

**Simplified from original spec:**
1. ❌ No `.shadowrc` config file - automatic tracking on first save
2. ✅ Single `list.json` metadata file instead of separate metadata files
3. ✅ Flat snapshot storage (`.shadow/snapshots/{id}`) - no directory hierarchy
4. ✅ Global config at `~/.config/sh_adow/config.yml` for flexible repo paths

**Core Concepts:**
- **Virtual HEAD**: Current file state (not stored, just referenced)
- **Version ID**: 8-char hex from SHA256 hash (e.g., `abc123ef`)
- **Auto-initialization**: `.shadow/` directory created on first save
- **Metadata**: All version info stored in single `list.json` file

### Storage Structure

```
.shadow/
├── list.json              # All metadata
└── snapshots/
    ├── abc123             # Flat file storage by version ID
    ├── def456
    └── xyz789
```

**list.json format:**
```json
{
  "files": [
    {
      "path": "/absolute/path/to/file.txt",
      "versions": [
        {
          "id": "abc123",
          "created_at": "2026-02-11T20:07:30Z",
          "tags": ["stable", "v1.0"],
          "notes": "Production config",
          "size": 1024,
          "hash": "sha256-full-hash..."
        }
      ]
    }
  ]
}
```

### Configuration

**Global config:** `~/.config/sh_adow/config.yml`

```yaml
repo_path: "./"  # Default: same directory as file
```

**Supported repo_path values:**
- `"./"` - Local (creates `.shadow/` in file's directory)
- `"~/.local/cache/"` - Absolute path (centralized storage)
- `"../cache/"` - Relative to file's directory

### Command Features

**save:**
- Auto-creates `.shadow/` if needed
- Interactive prompts for tags/notes (optional)
- Flags: `-t/--tag`, `-n/--note`

**list:**
- Shows all tracked files (no args)
- Shows specific file versions (with file arg)
- Displays: version ID, age, tags, notes, size

**restore:**
- Interactive prompt to save current state before restoring
- Flag: `--no-save` to skip prompt

**delete:**
- Interactive confirmation prompt
- Flag: `-f/--force` to skip confirmation

### Tech Stack

Built with Charm ecosystem for beautiful CLI:
- **cobra** - Command structure and routing
- **huh** - Interactive forms and prompts
- **lipgloss** - Terminal styling and colors
- **yaml.v3** - Config file parsing

### Project Structure

```
sh_adow/
├── cmd/                   # Cobra commands
│   ├── root.go           # CLI root + setup
│   ├── save.go           # Save command
│   ├── list.go           # List command
│   ├── restore.go        # Restore command
│   └── delete.go         # Delete command
├── internal/
│   ├── config/           # Config parser (~/.config/sh_adow/config.yml)
│   │   └── config.go
│   ├── repo/             # Shadow path resolution
│   │   └── repo.go
│   └── shadow/           # Core versioning logic
│       └── list.go       # Metadata and file operations
├── main.go               # Entry point
├── go.mod                # Dependencies
└── README.md             # User documentation
```

### Design Principles Followed

✅ **Simplicity** - No complex git concepts, no branching/merging  
✅ **Lightweight** - ~6.7MB binary, minimal dependencies  
✅ **User-friendly** - Interactive prompts with fallback flags  
✅ **Flexible** - Configurable storage paths  
✅ **Zero setup** - Works immediately, auto-initializes  

### Known Limitations

- No diff viewing (use external tools)
- No remote sync or backup
- No compression (simplicity over space)
- No directory versioning (files only)
- No file watching/auto-save

### Future Enhancements (Optional)

- `shadow diff <file> <version-id>` - Show differences
- `shadow export <file> <version-id>` - Export version to new file
- `shadow gc` - Garbage collect orphaned snapshots
- Compression option in config
- Shell completion scripts

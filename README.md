Notion App
==========

About
-----
Notion App is a terminal (CLI) note-taking application written in Go. It creates, updates, and lists Markdown files inside a local vault directory. The UI is implemented using the Bubble Tea ecosystem (bubbles, textarea, list, textinput, and lipgloss) for a polished, keyboard-driven experience.

Key features
------------
- Create new Markdown (.md) notes.
- Edit and save notes.
- List existing notes with simple metadata (modified time).

What I used
-----------
- Language: Go
- UI: Bubble Tea ecosystem — github.com/charmbracelet/bubbletea, bubbles/textinput, bubbles/textarea, bubbles/list, and lipgloss.
- Project layout: single executable with clear sections in main.go (init, model, updates, views currently in one file).

How to run
----------
Prerequisites:
- Go 1.24+

Run with Make (recommended):

```sh
make run
```

Or run directly from the project root:

```sh
go run main.go
```

Operation / Working
-------------------
The app follows this runtime flow:

1. Initialization: the vault directory (e.g. ~/.notion in the current working directory) is created if missing.
2. UI bootstrap: textinput, textarea, and list models are initialized for user interaction.
3. User actions:
   - Ctrl+N opens the new-file input to create a note (saved as filename.md).
   - Ctrl+L lists notes in the vault; pressing enter opens a selected note for editing.
   - Ctrl+S saves the current note to disk.
   - Esc navigates back / clears current state; Ctrl+Q quits.

Files of interest
-----------------
- [main.go](main.go) — primary executable and the app logic (UI, handlers, and file operations).
- [init/init.go](init/init.go) — initialization helpers (where present).
- [model/model.go](model/model.go) — data structures and model helpers.
- [updates/update.go](updates/update.go) — business logic for updates.
- [views/view.go](views/view.go) — rendering helpers.

Improvements / Roadmap
---------------------
- Split main.go into focused packages: init, model, updates, and views (move UI handling and business logic into separate files).
- Implement a delete-file operation (keyboard shortcut and confirmation) to remove notes from the vault.
- Add persistent metadata (DB) or index for faster search and more metadata.
- Add unit and integration tests for file operations and update logic.
- Add CLI flags and a config file for customizing vault location and settings.
- Consider exposing a minimal HTTP API or web UI for remote access and syncing.

Next steps I can help with
-------------------------
- Implement the delete-file feature and its keyboard binding.
- Refactor main.go into views, updates, model, and init files.
- Add tests and a Makefile target for lint/test/build.

If you want me to start any of the improvements above, tell me which one and I'll implement it.

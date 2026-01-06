# Task Tracker CLI

Simple command-line task tracker that stores tasks in a local JSON file. Supports adding, updating, deleting, listing, and marking tasks as in-progress or done using positional arguments only—no external dependencies.

## Features
- Add, update, delete tasks.
- Mark tasks as `todo`, `in-progress`, or `done`.
- List all tasks or filter by status (`todo`, `in-progress`, `done`).
- Persists tasks to `db.json` in the current directory; file is created if missing.

## Installation
```bash
git clone https://github.com/Rexbrainz/task-tracker.git
cd task-tracker
go build -o task-cli
```

## Usage
All commands use positional arguments. Examples:
```bash
# Add a new task
./task-cli add "Buy groceries"

# Update and delete
./task-cli update 1 "Buy groceries and cook dinner"
./task-cli delete 1

# Mark status
./task-cli mark-in-progress 1
./task-cli mark-done 1

# List tasks
./task-cli list          # all
./task-cli list done     # only done
./task-cli list todo     # only todo
./task-cli list in-progress # only in-progress
```

## Task Data Model
Each task is stored in `db.json` with:
- `id` (int): unique identifier
- `description` (string): short description
- `status` (string): one of `todo`, `in-progress`, `done`
- `createdAt` (string, RFC3339)
- `updatedAt` (string, RFC3339)

`db.json` lives in the working directory where you run the CLI.

## Notes
- Uses only the Go standard library and native filesystem access.
- Handles basic validation (missing args, empty descriptions, unknown IDs, invalid status filters).

## Project Source
This project is one of the tasks from https://roadmap.sh/projects/task-tracker
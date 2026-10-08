# TaskI

TaskI is a small command-line task tracker written in Go. Create tasks, update their descriptions or status, and list them from your terminal. Tasks are saved as JSON in the directory where you run the command.

Inspired by the [roadmap.sh Task Tracker project](https://roadmap.sh/projects/task-tracker).

## Features

- Create, update, and delete tasks
- Track tasks as `todo`, `in progress`, or `done`
- List every task or filter by status
- Persist tasks in a local `tasks.json` file

## Requirements

- Go 1.26.5 or later

## Install

Clone the repository and install the `TaskI` executable:

```sh
git clone https://github.com/heisAnselem/TaskI.git
cd TaskI
go install .
```

Make sure your Go `bin` directory is on your `PATH`. You can also run the application from the repository with `go run .`.

## Usage

The commands below use `TaskI` as the executable name. If you run from the repository with `go run .`, replace `TaskI` with `go run .`.

### Create a task

```sh
TaskI create -description "Buy coffee"
```

New tasks start with the `todo` status.

### List tasks

```sh
# List all tasks
TaskI list

# List tasks with a specific status
TaskI list -status "done"
TaskI list -status "in progress"
TaskI list -status "todo"
```

### Update a task

Provide a task ID and at least one field to change:

```sh
# Change the status
TaskI update -id 1 -status "in progress"

# Change the description
TaskI update -id 1 -description "Buy tea"

# Change both
TaskI update -id 1 -status "done" -description "Buy tea"
```

Valid statuses are `todo`, `in progress`, and `done`.

### Delete a task

```sh
TaskI delete -id 1
```

## Data storage

TaskI reads and writes `tasks.json` in the current working directory. The file is created when you add the first task, so running TaskI from different directories keeps those task lists separate.

Each task stores an ID, description, status, creation time, and last updated time. IDs are assigned automatically.

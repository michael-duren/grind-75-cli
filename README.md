# Grind 75 CLI

## About

Grind 75 CLI is a CLI tool designed to present a list of leet code questions from the popular Grind 75 coding challenge.
It allows users to keep track of their progress in their terminal and select their preferred plan (number of weeks and hours per week) which will update the list of questions accordingly.

## Design

### TUI

G7C is built with the charm CLI libraries:

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - TUI styles
- [Harmonica](https://github.com/charmbracelet/harmonica) - For fun animations
- [Bubbles](https://github.com/charmbracelet/bubbles) - Reusable TUI components
- [Log](https://github.com/charmbracelet/log) - My favorite logging library

I know go devs like to do everything themselves, but when it comes to any UI I'm lazy.

### Data

G7C will save your progress in a sqlite database located at `~/.g7c/g7c.db`.
The database is managed using [sqlc](https://sqlc.dev/) to generate type-safe queries.
Otherwise it's just raw SQL.

You can optionally import different problems or customize your own grind plan by configuring
a `config.yaml` file located at `~/.g7c/config.yaml` (more information on the `config.yaml` can be found below).

Probles must be in a JSON file in the following format located at `~/.g7c/problems.json`:

```json
[
  {
    "slug": "two-sum",
    "title": "Two Sum",
    "url": "https://leetcode.com/problems/two-sum/",
    "duration": 15,
    "epi": null,
    "difficulty": "Easy",
    "id": 1,
    "topic": "array",
    "routines": ["hashing"]
  }
]
```

### Navigation

G7C uses keyboard navigation to move around the TUI, there are three main views:

- **Main View** - View and track your progress on the grind questions
- **Settings View** - Select your grind plan (weeks and hours per week), optionally select
  custom problem set. You may also setup reminders via Email if you've configured SMTP in
  the config file. There are also export options here to export your progress to CSV or JSON.
- **Help View** - View the keybindings and commands available

## Features

- View and track your progress on the Grind 75 questions
- Select your grind plan (weeks and hours per week)
- Optionally import custom problem sets via JSON file
- Export your progress to CSV or JSON
- Setup email reminders for daily/weekly progress (requires SMTP configuration)
- Select struggle problems to later review

# Grind 75 CLI (G7C)

A terminal-based progress tracker for the Grind 75 coding challenge with customizable study plans and email reminders.

## Table of Contents

- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Commands](#commands)
  - [TUI Navigation](#tui-navigation)
- [Configuration](#configuration)
  - [Default Example](#default-example)
  - [Full Example](#full-example)
  - [Custom Problem Sets](#custom-problem-sets)
- [Design](#design)
  - [TUI](#tui)
  - [Data Storage](#data-storage)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## About

Grind 75 CLI is a CLI tool designed to present a list of leet code questions from the popular Grind 75 coding challenge.
It allows users to keep track of their progress in their terminal and select their preferred plan (number of weeks and hours per week) which will update the list of questions accordingly.

## Features

- View and track your progress on the Grind 75 questions. This includes fields for time taken, date finished, review date (if applicable), number of attempts to solve and status (completed, uncompleted, struggling).
- Select your grind plan (weeks and hours per week)
- Optionally import custom problem sets via JSON file
- Export your progress to CSV or JSON
- Setup email reminders for daily/weekly progress (requires SMTP configuration)
- Select problems you've struggled with for additional practice, and view only
  those problems. Set a reveiw date to remind you to revisit them.

## Installation

You can install Grind 75 CLI by downloading the latest release from the [releases page](https://github.com/michael-duren/grind-75-cli/releases).
Alternatively, if you have Go installed, you can run:

```bash
go install github.com/michael-duren/grind-75-cli@latest
```

After installation, you can run the CLI by executing:

```bash
g7c
```

## Quick Start

1. Install: `go install github.com/michael-duren/grind-75-cli@latest`
2. Run: `g7c`
3. (Optional) Setup email reminders: `g7c setup-email` then `g7c cron on`

## Usage

### Commands

G7C is built using [Cobra](https://github.com/spf13/cobra) for command line parsing.
The main commands are:

- `g7c` - Launch the TUI
- `g7c setup-email` - Setup SMTP email configuration (stores password in system keychain)
- `g7c check-reminders` - Check for any reminders to send (can be added to cron)
- `g7c cron`
  - `on` - Setup a cron job to run `g7c check-reminders` every day at midnight
  - `off` - Remove the cron job

### TUI Navigation

G7C uses keyboard navigation to move around the TUI, there are three main views:

- **Main View** - View and track your progress on the grind questions. You can mark problems as
  completed, view problem details, and filter problems by status (all, completed, uncompleted, struggling).
- **Settings View** - Select your grind plan (weeks and hours per week), optionally select
  custom problem set. You may also setup reminders via Email if you've configured SMTP in
  the config file. There are also export options here to export your progress to CSV or JSON.
- **Help View** - View the keybindings and commands available (reachable from any view by pressing `?`).

## Configuration

[!IMPORTANT]
**A configuration file is optional. A default file will be created at `~/.g7c/config.yaml` on first run if it does not exist.**

Grind 75 CLI can be configured via a `config.yaml` file located at `~/.g7c/config.yaml`.
Here is an example configuration file:

### Full Example

```yaml
grind_plan:
  weeks: 10
  hours_per_week: 5
smtp:
  provider: "gmail"
  email: "myemail@gmail.com"
  # WARNING: Storing password here is NOT SECURE
  # RECOMMENDED: Use G7C_SMTP_PASSWORD env var or run 'g7c setup-email' to store in system keychain
  password: "" # Leave empty and use env var or keychain instead
# Reminders require SMTP configuration and a background scheduler
# Run 'g7c cron on' to enable automatic reminders
reminders:
  daily_reminder: true
  daily_time: "09:00"
  weekly_overview: false
  weekly_day: "monday"
  weekly_time: "09:00"
  struggle_reminder: true
  struggle_days_before: 1 # Remind 1 day before review date
custom_problem_set:
  enabled: false
  path: "~/.g7c/problems.json"
  # These override grind_plan settings when enabled
  weeks: 12
  hours_per_week: 6
```

### Default Example

```yaml
grind_plan:
  weeks: 10
  hours_per_week: 5
smtp:
  provider: ""
  email: ""
  password: ""
reminders:
  daily_reminder: false
  daily_time: "09:00"
  weekly_overview: false
  weekly_day: "monday"
  weekly_time: "09:00"
  struggle_reminder: false
  struggle_days_before: 1
custom_problem_set:
  enabled: false
```

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

## Troubleshooting

**Email reminders not sending?**

- Check `G7C_SMTP_PASSWORD` is set or run `g7c setup-email`
- Verify cron is running: `g7c cron on`
- Test manually: `g7c check-reminders`

**Gmail not working?**

- Enable 2FA on your Google account
- Generate an App Password: [instructions](https://support.google.com/accounts/answer/185833)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

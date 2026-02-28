# tomato-cli

A chat-based Pomodoro CLI tool for focused work sessions, inspired by Claude Code's interactive interface.

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or later installed

### Linux/macOS

```bash
# Clone the repository
git clone https://github.com/quoder/tomato-cli.git
cd tomato-cli

# Build the CLI
make build

# Run it
./bin/tomato
```

Or install globally:

```bash
sudo cp bin/tomato /usr/local/bin/
tomato
```

### Windows

```cmd
# Clone the repository
git clone https://github.com/quoder/tomato-cli.git
cd tomato-cli

# Build the CLI
go build -o bin\tomato.exe .\cmd\tomato

# Run it
bin\tomato.exe
```

Or install to your PATH:

```cmd
copy bin\tomato.exe %USERPROFILE%\AppData\Local\Programs\tomato\
setx PATH "%PATH%;%USERPROFILE%\AppData\Local\Programs\tomato"
```

## Usage

Run the CLI:

```bash
./bin/tomato
```

You'll see a chat prompt:

```
🍅 tomato-cli - Pomodoro Timer
Type /help for available commands

> 
```

### Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `/start [work] [break]` | - | Start pomodoro (default: 25m work, 5m break) |
| `/pause` | - | Pause current timer |
| `/resume` | - | Resume paused timer |
| `/next` | - | Skip to next phase (cycles continuously) |
| `/cancel` | `/stop` | Cancel current session |
| `/status` | - | Show current status |
| `/stats` | - | Show completed pomodoros |
| `/config` | - | Show configuration |
| `/help` | - | Show help |
| `/exit` | `/quit` | Exit CLI |

### Examples

```bash
# Start with default 25min work, 5min break
> /start

# Start with custom durations (10 min work, 3 min break)
> /start 10 3

# Check status
> /status

# Pause and resume
> /pause
> /resume

# Skip to next phase (cycles: work → break → work → break...)
> /next

# Cancel/stop the session
> /cancel
# or
> /stop

# View statistics
> /stats

# Exit
> /exit
# or
> /quit
```

## Features

- **Continuous cycling**: `/next` cycles through work and break phases indefinitely
- **Custom durations**: Durations specified in `/start` persist across all cycles
- **Smart start**: If timer is paused, `/start` shows a warning instead of overwriting

## Storage

Configuration and statistics are stored in:

- **Linux/macOS**: `~/.config/tomato-cli/`
- **Windows**: `%APPDATA%\tomato-cli\`

Files:
- `config.json` - User preferences
- `stats.json` - Completed session history

## Project Structure

```
tomato-cli/
├── cmd/
│   └── tomato/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go       # JSON config handling
│   ├── timer/
│   │   └── timer.go        # Channel-driven timer
│   ├── session/
│   │   └── session.go      # Session orchestration
│   ├── stats/
│   │   └── stats.go       # Statistics tracking
│   └── repl/
│       └── repl.go         # REPL interface
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Go Concepts Used

1. **Goroutines & Channels** - Timer running in background
2. **JSON encoding/decoding** - Config & stats persistence
3. **File I/O** config files
4 - Reading/writing. **Command pattern** - REPL command handling
5. **State machine** - Session phases

package repl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/quoder/tomato-cli/internal/session"
)

type REPL struct {
	scanner *bufio.Scanner
	session *session.Session
	running bool
}

func New(sess *session.Session) *REPL {
	return &REPL{
		scanner: bufio.NewScanner(os.Stdin),
		session: sess,
		running: true,
	}
}

func (r *REPL) Run() {
	fmt.Println("🍅 tomato-cli - Pomodoro Timer")
	fmt.Println("Type /help for available commands")
	fmt.Println()

	for r.running {
		fmt.Print("> ")
		if !r.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.scanner.Text())
		if line == "" {
			continue
		}

		r.handleCommand(line)
	}
}

func (r *REPL) handleCommand(line string) {
	if !strings.HasPrefix(line, "/") {
		fmt.Println("Unknown command. Type /help for available commands")
		return
	}

	parts := strings.Fields(line)
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/start":
		r.handleStart(args)
	case "/pause":
		r.handlePause()
	case "/resume":
		r.handleResume()
	case "/next":
		r.handleNext()
	case "/cancel":
		r.handleCancel()
	case "/stop":
		r.handleCancel()
	case "/status":
		r.handleStatus()
	case "/stats":
		r.handleStats()
	case "/config":
		r.handleConfig()
	case "/help":
		r.handleHelp()
	case "/exit":
		r.handleEnd()
	case "/quit":
		r.handleEnd()
	default:
		fmt.Printf("Unknown command: %s. Type /help for available commands\n", cmd)
	}
}

func (r *REPL) handleStart(args []string) {
	workMins := 0
	breakMins := 0

	if len(args) >= 1 {
		if v, err := strconv.Atoi(args[0]); err == nil {
			workMins = v
		}
	}
	if len(args) >= 2 {
		if v, err := strconv.Atoi(args[1]); err == nil {
			breakMins = v
		}
	}

	if err := r.session.Start(workMins, breakMins); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if workMins > 0 && breakMins > 0 {
		fmt.Printf("Started: %d min work, %d min break\n", workMins, breakMins)
	} else {
		fmt.Println("Started: 25 min work, 5 min break")
	}
}

func (r *REPL) handlePause() {
	if err := r.session.Pause(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Timer paused")
}

func (r *REPL) handleResume() {
	if err := r.session.Resume(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Timer resumed")
}

func (r *REPL) handleNext() {
	if err := r.session.Next(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Skipped to next phase")
}

func (r *REPL) handleCancel() {
	if err := r.session.Cancel(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Session cancelled")
}

func (r *REPL) handleStatus() {
	status := r.session.GetStatus()
	fmt.Printf("Status: %s\n", status)
}

func (r *REPL) handleStats() {
	total := r.session.GetTotalCompleted()
	today := r.session.GetTodayCompleted()
	fmt.Printf("Statistics:\n")
	fmt.Printf("  Total completed: %d\n", total)
	fmt.Printf("  Today completed: %d\n", today)
}

func (r *REPL) handleConfig() {
	cfg := r.session.GetConfig()
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Work duration: %v\n", cfg.WorkDuration)
	fmt.Printf("  Break duration: %v\n", cfg.BreakDuration)
}

func (r *REPL) handleHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  /start [work] [break]  - Start pomodoro (e.g., /start 25 5)")
	fmt.Println("  /pause                 - Pause current timer")
	fmt.Println("  /resume                - Resume paused timer")
	fmt.Println("  /next                  - Skip to next phase")
	fmt.Println("  /cancel /stop          - Cancel current session")
	fmt.Println("  /status                - Show current status")
	fmt.Println("  /stats                 - Show statistics")
	fmt.Println("  /config                - Show/edit configuration")
	fmt.Println("  /help                  - Show this help")
	fmt.Println("  /exit /quit             - Exit")
}

func (r *REPL) handleEnd() {
	r.running = false
	fmt.Println("Goodbye!")
}

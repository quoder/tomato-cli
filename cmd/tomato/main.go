package main

import (
	"github.com/quoder/tomato-cli/internal/config"
	"github.com/quoder/tomato-cli/internal/repl"
	"github.com/quoder/tomato-cli/internal/session"
	"github.com/quoder/tomato-cli/internal/stats"
)

func main() {
	cfg, err := config.NewManager()
	if err != nil {
		println("Error loading config:", err)
		return
	}

	st, err := stats.NewManager()
	if err != nil {
		println("Error loading stats:", err)
		return
	}

	sess := session.New(cfg, st)

	r := repl.New(sess)
	r.Run()
}

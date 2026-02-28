package timer

import (
	"time"
)

type Phase string

const (
	PhaseIdle        Phase = "idle"
	PhaseWork        Phase = "work"
	PhaseWorkPaused  Phase = "work_paused"
	PhaseBreak       Phase = "break"
	PhaseBreakPaused Phase = "break_paused"
)

type State struct {
	Phase         Phase
	Remaining     time.Duration
	TotalDuration time.Duration
}

type command struct {
	action   string
	duration time.Duration
	phase    Phase
	respChan chan State
}

type Timer struct {
	cmdChan chan command
	state   State
}

func New() *Timer {
	t := &Timer{
		cmdChan: make(chan command),
		state:   State{Phase: PhaseIdle},
	}
	go t.run()
	return t
}

func (t *Timer) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-t.cmdChan:
			t.handleCommand(cmd)
		case <-ticker.C:
			t.tick()
		}
	}
}

func (t *Timer) handleCommand(cmd command) {
	switch cmd.action {
	case "start":
		t.state = State{
			Phase:         cmd.phase,
			Remaining:     cmd.duration,
			TotalDuration: cmd.duration,
		}
	case "pause":
		t.state.Phase = pausedPhase(t.state.Phase)
	case "resume":
		t.state.Phase = resumedPhase(t.state.Phase)
	case "stop":
		t.state = State{Phase: PhaseIdle, Remaining: 0}
	case "get":
		cmd.respChan <- t.state
	}
}

func (t *Timer) tick() {
	if t.state.Phase == PhaseIdle || IsPausedPhase(t.state.Phase) {
		return
	}

	t.state.Remaining -= time.Second
	if t.state.Remaining <= 0 {
		t.state = State{Phase: PhaseIdle, Remaining: 0}
	}
}

func (t *Timer) Start(duration time.Duration, phase Phase) {
	t.cmdChan <- command{action: "start", duration: duration, phase: phase}
}

func (t *Timer) Pause() {
	t.cmdChan <- command{action: "pause"}
}

func (t *Timer) Resume() {
	t.cmdChan <- command{action: "resume"}
}

func (t *Timer) Stop() {
	t.cmdChan <- command{action: "stop"}
}

func (t *Timer) GetState() State {
	respChan := make(chan State, 1)
	t.cmdChan <- command{action: "get", respChan: respChan}
	return <-respChan
}

func (t *Timer) IsRunning() bool {
	state := t.GetState()
	return state.Phase != PhaseIdle
}

func (t *Timer) IsPaused() bool {
	state := t.GetState()
	return IsPausedPhase(state.Phase)
}

func pausedPhase(phase Phase) Phase {
	switch phase {
	case PhaseWork:
		return PhaseWorkPaused
	case PhaseBreak:
		return PhaseBreakPaused
	default:
		return phase
	}
}

func resumedPhase(phase Phase) Phase {
	switch phase {
	case PhaseWorkPaused:
		return PhaseWork
	case PhaseBreakPaused:
		return PhaseBreak
	default:
		return phase
	}
}

func ActivePhase(phase Phase) Phase {
	switch phase {
	case PhaseWorkPaused:
		return PhaseWork
	case PhaseBreakPaused:
		return PhaseBreak
	default:
		return phase
	}
}

func IsPausedPhase(phase Phase) bool {
	return phase == PhaseWorkPaused || phase == PhaseBreakPaused
}

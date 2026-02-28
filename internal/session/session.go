package session

import (
	"errors"
	"time"

	"github.com/quoder/tomato-cli/internal/config"
	"github.com/quoder/tomato-cli/internal/stats"
	"github.com/quoder/tomato-cli/internal/timer"
)

var ErrTimerPaused = errors.New("tomato is paused, please /resume or /cancel to start a new one")

type Session struct {
	config        *config.Manager
	stats         *stats.Manager
	timer         *timer.Timer
	startTime     time.Time
	workDuration  time.Duration
	breakDuration time.Duration
}

func New(cfg *config.Manager, st *stats.Manager) *Session {
	s := &Session{
		config: cfg,
		stats:  st,
		timer:  timer.New(),
	}
	s.timer.SetCallback(s.onPhaseComplete)
	return s
}

func (s *Session) onPhaseComplete(completedState timer.State) {
	activePhase := timer.ActivePhase(completedState.Phase)
	if completedState.TotalDuration > 0 {
		s.recordSession(activePhase, completedState.TotalDuration)
	}
	if activePhase == timer.PhaseWork {
		s.startTime = time.Now()
		s.timer.Start(s.breakDuration, timer.PhaseBreak)
	} else if activePhase == timer.PhaseBreak {
		s.startTime = time.Now()
		s.timer.Start(s.workDuration, timer.PhaseWork)
	}
}

func (s *Session) Start(workMinutes, breakMinutes int) error {
	if s.timer.IsRunning() && s.timer.IsPaused() {
		return ErrTimerPaused
	}
	if s.timer.IsRunning() {
		return nil
	}

	cfg := s.config.Get()
	workDuration := cfg.WorkDuration

	if workMinutes > 0 {
		workDuration = time.Duration(workMinutes) * time.Minute
	}

	breakDuration := cfg.BreakDuration
	if breakMinutes > 0 {
		breakDuration = time.Duration(breakMinutes) * time.Minute
	}

	s.workDuration = workDuration
	s.breakDuration = breakDuration

	s.startTime = time.Now()

	s.timer.Start(workDuration, timer.PhaseWork)

	return nil
}

func (s *Session) Pause() error {
	if !s.timer.IsRunning() {
		return nil
	}
	s.timer.Pause()
	return nil
}

func (s *Session) Resume() error {
	s.timer.Resume()
	return nil
}

func (s *Session) Next() error {
	if !s.timer.IsRunning() {
		return nil
	}

	currentState := s.timer.GetState()
	phase := timer.ActivePhase(currentState.Phase)

	switch phase {
	case timer.PhaseWork:
		s.recordSession(timer.PhaseWork, currentState.TotalDuration)
		s.timer.Start(s.breakDuration, timer.PhaseBreak)
	case timer.PhaseBreak:
		s.timer.Start(s.workDuration, timer.PhaseWork)
	}

	return nil
}

func (s *Session) Cancel() error {
	if !s.timer.IsRunning() {
		return nil
	}

	currentState := s.timer.GetState()
	s.timer.Stop()
	s.recordSession(timer.ActivePhase(currentState.Phase), currentState.TotalDuration-currentState.Remaining)
	return nil
}

func (s *Session) GetStatus() string {
	if !s.timer.IsRunning() {
		return "Idle"
	}

	state := s.timer.GetState()
	phase := "Work"
	if timer.ActivePhase(state.Phase) == timer.PhaseBreak {
		phase = "Break"
	}

	paused := ""
	if timer.IsPausedPhase(state.Phase) {
		paused = " (Paused)"
	}

	return phase + ": " + state.Remaining.String() + paused
}

func (s *Session) GetState() timer.State {
	return s.timer.GetState()
}

func (s *Session) IsRunning() bool {
	return s.timer.IsRunning()
}

func (s *Session) recordSession(phase timer.Phase, duration time.Duration) {
	session := stats.Session{
		StartTime: s.startTime,
		EndTime:   time.Now(),
		Duration:  duration,
		Phase:     string(phase),
	}
	s.stats.AddSession(session)
}

func (s *Session) GetTotalCompleted() int {
	return s.stats.GetTotalCompleted()
}

func (s *Session) GetTodayCompleted() int {
	return s.stats.GetTodayCompleted()
}

func (s *Session) ClearStats() error {
	return s.stats.Clear()
}

func (s *Session) GetConfig() config.Config {
	return s.config.Get()
}

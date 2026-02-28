package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/quoder/tomato-cli/internal/storage"
)

type Session struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Phase     string        `json:"phase"`
}

type Stats struct {
	Sessions []Session `json:"sessions"`
}

type Manager struct {
	statsPath string
	stats     Stats
}

func NewManager() (*Manager, error) {
	configDir, err := storage.AppDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	statsPath := filepath.Join(configDir, "stats.json")
	m := &Manager{
		statsPath: statsPath,
		stats:     Stats{Sessions: []Session{}},
	}

	if err := m.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) load() error {
	data, err := os.ReadFile(m.statsPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.stats)
}

func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.statsPath, data, 0644)
}

func (m *Manager) AddSession(session Session) error {
	m.stats.Sessions = append(m.stats.Sessions, session)
	return m.save()
}

func (m *Manager) GetSessions() []Session {
	return m.stats.Sessions
}

func (m *Manager) GetTotalCompleted() int {
	count := 0
	for _, s := range m.stats.Sessions {
		if s.Phase == "work" {
			count++
		}
	}
	return count
}

func (m *Manager) GetTodayCompleted() int {
	count := 0
	today := time.Now().Truncate(24 * time.Hour)
	for _, s := range m.stats.Sessions {
		if s.Phase == "work" && s.StartTime.After(today) {
			count++
		}
	}
	return count
}

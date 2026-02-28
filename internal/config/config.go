package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/quoder/tomato-cli/internal/storage"
)

type Config struct {
	WorkDuration  time.Duration `json:"work_duration"`
	BreakDuration time.Duration `json:"break_duration"`
}

var Default = Config{
	WorkDuration:  25 * time.Minute,
	BreakDuration: 5 * time.Minute,
}

type Manager struct {
	configPath string
	config     Config
}

func NewManager() (*Manager, error) {
	configDir, err := storage.AppDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")
	m := &Manager{
		configPath: configPath,
		config:     Default,
	}

	if err := m.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Save the default config.json in first run
		if err := m.save(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Manager) load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.config)
}

func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath, data, 0644)
}

func (m *Manager) Get() Config {
	return m.config
}

func (m *Manager) SetWorkDuration(d time.Duration) error {
	m.config.WorkDuration = d
	return m.save()
}

func (m *Manager) SetBreakDuration(d time.Duration) error {
	m.config.BreakDuration = d
	return m.save()
}

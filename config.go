package main

import (
	"encoding/json"
	"os"
)

// Config disimpan di file config.json di folder yang sama dengan .exe
// Supaya token tidak hardcode di dalam binary.
type Config struct {
	ServerURL       string `json:"server_url"`        // contoh: https://print.marsa9.com
	AgentToken      string `json:"agent_token"`       // didapat saat register manual di dashboard
	PollIntervalSec int    `json:"poll_interval_sec"` // default 20-30 detik
	HeartbeatSec    int    `json:"heartbeat_sec"`     // default 60 detik
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// Fallback default kalau nilai kosong/nol
	if cfg.PollIntervalSec == 0 {
		cfg.PollIntervalSec = 20
	}
	if cfg.HeartbeatSec == 0 {
		cfg.HeartbeatSec = 60
	}
	return &cfg, nil
}

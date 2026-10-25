package main

import (
	"sync"
	"time"
)

// JobEvent satu baris riwayat job untuk ditampilkan di dashboard
type JobEvent struct {
	JobID     int       `json:"job_id"`
	Printer   string    `json:"printer"`
	Status    string    `json:"status"` // printing / done / failed
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// RuntimeState nyimpan kondisi Agent saat ini (thread-safe), dipakai
// baik oleh loop polling/heartbeat maupun web dashboard secara bersamaan
type RuntimeState struct {
	mu            sync.RWMutex
	lastHeartbeat time.Time
	lastPoll      time.Time
	startedAt     time.Time
	recentJobs    []JobEvent
	maxJobs       int
}

func NewRuntimeState() *RuntimeState {
	return &RuntimeState{
		startedAt: time.Now(),
		maxJobs:   30,
	}
}

func (r *RuntimeState) MarkHeartbeat() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastHeartbeat = time.Now()
}

func (r *RuntimeState) MarkPoll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastPoll = time.Now()
}

func (r *RuntimeState) AddJobEvent(ev JobEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ev.Timestamp = time.Now()
	r.recentJobs = append([]JobEvent{ev}, r.recentJobs...) // terbaru di atas
	if len(r.recentJobs) > r.maxJobs {
		r.recentJobs = r.recentJobs[:r.maxJobs]
	}
}

// Snapshot data buat ditampilkan di dashboard (copy, bukan referensi langsung)
type Snapshot struct {
	LastHeartbeat time.Time
	LastPoll      time.Time
	StartedAt     time.Time
	RecentJobs    []JobEvent
}

func (r *RuntimeState) Snapshot() Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	jobsCopy := make([]JobEvent, len(r.recentJobs))
	copy(jobsCopy, r.recentJobs)
	return Snapshot{
		LastHeartbeat: r.lastHeartbeat,
		LastPoll:      r.lastPoll,
		StartedAt:     r.startedAt,
		RecentJobs:    jobsCopy,
	}
}

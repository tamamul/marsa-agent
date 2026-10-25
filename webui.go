package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const dashboardHTML = `
<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>MARSA Print Agent</title>
<style>
  :root {
    --accent: #0078D4;
    --bg: #F3F3F3;
    --card-bg: rgba(255,255,255,0.7);
    --text: #1A1A1A;
    --text-dim: #5c5c5c;
    --border: rgba(0,0,0,0.07);
    --shadow: 0 1px 3px rgba(0,0,0,0.08), 0 4px 12px rgba(0,0,0,0.04);
  }
  @media (prefers-color-scheme: dark) {
    :root {
      --bg: #202020;
      --card-bg: rgba(255,255,255,0.05);
      --text: #F3F3F3;
      --text-dim: #b0b0b0;
      --border: rgba(255,255,255,0.08);
      --shadow: 0 1px 3px rgba(0,0,0,0.3), 0 4px 16px rgba(0,0,0,0.2);
    }
  }
  * { box-sizing: border-box; }
  body {
    font-family: "Segoe UI Variable", "Segoe UI", system-ui, sans-serif;
    max-width: 760px; margin: 0 auto; padding: 32px 20px 60px;
    color: var(--text);
    background:
      radial-gradient(1200px 500px at 10% -10%, rgba(0,120,212,0.12), transparent),
      var(--bg);
    min-height: 100vh;
  }
  h1 { font-size: 22px; font-weight: 600; margin: 0 0 20px; display: flex; align-items: center; gap: 10px; }
  .card {
    background: var(--card-bg);
    backdrop-filter: blur(20px);
    border: 1px solid var(--border);
    border-radius: 12px;
    box-shadow: var(--shadow);
    padding: 18px 20px;
    margin-bottom: 16px;
  }
  .card h2 { font-size: 14px; font-weight: 600; margin: 0 0 12px; color: var(--text-dim); text-transform: uppercase; letter-spacing: .03em; }
  .ok { color: #0F9D58; font-weight: 600; }
  .warn { color: #D13438; font-weight: 600; }
  .meta { color: var(--text-dim); font-size: 13px; margin-top: 4px; }
  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  td, th { text-align: left; padding: 8px 6px; border-bottom: 1px solid var(--border); }
  th { color: var(--text-dim); font-weight: 600; font-size: 12px; }
  .status-done { color: #0F9D58; font-weight: 600; }
  .status-failed { color: #D13438; font-weight: 600; }
  .status-printing { color: #CA5010; font-weight: 600; }
  label { display: block; font-size: 12px; color: var(--text-dim); margin: 10px 0 4px; }
  input {
    width: 100%; padding: 9px 10px; border-radius: 6px;
    border: 1px solid var(--border); background: var(--bg); color: var(--text);
    font-size: 14px; outline: none;
  }
  input:focus { border-color: var(--accent); box-shadow: 0 0 0 2px rgba(0,120,212,0.25); }
  button {
    margin-top: 16px; padding: 9px 20px; background: var(--accent); color: white;
    border: none; border-radius: 6px; cursor: pointer; font-size: 14px; font-weight: 500;
  }
  button:hover { background: #106EBE; }
  .msg { padding: 10px 14px; border-radius: 8px; margin-bottom: 16px; font-size: 14px; }
  .msg-ok { background: rgba(15,157,88,0.12); color: #0F9D58; border: 1px solid rgba(15,157,88,0.25); }
  .dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; margin-right: 6px; }
  .dot-ok { background: #0F9D58; }
  .dot-warn { background: #D13438; }
</style>
</head>
<body>
<h1>🖨️ MARSA Print Agent</h1>

{{if .Saved}}<div class="msg msg-ok">Setting tersimpan. Diterapkan langsung tanpa perlu restart.</div>{{end}}

<div class="card">
  <h2>Status</h2>
  {{if .IsHealthy}}<span class="dot dot-ok"></span><span class="ok">Berjalan normal</span>{{else}}<span class="dot dot-warn"></span><span class="warn">Belum ada heartbeat terbaru</span>{{end}}
  <div class="meta">Jalan sejak {{.StartedAt}}</div>
  <div class="meta">Heartbeat terakhir: {{.LastHeartbeat}}</div>
  <div class="meta">Polling terakhir: {{.LastPoll}}</div>
</div>

<div class="card">
  <h2>Riwayat Job Terakhir</h2>
  <table>
    <tr><th>Waktu</th><th>Job</th><th>Printer</th><th>Status</th><th>Keterangan</th></tr>
    {{range .RecentJobs}}
    <tr>
      <td>{{.Timestamp.Format "15:04:05"}}</td>
      <td>#{{.JobID}}</td>
      <td>{{.Printer}}</td>
      <td class="status-{{.Status}}">{{.Status}}</td>
      <td>{{.Message}}</td>
    </tr>
    {{else}}
    <tr><td colspan="5">Belum ada job masuk</td></tr>
    {{end}}
  </table>
</div>

<div class="card">
  <h2>Setting</h2>
  <form method="POST" action="/settings">
    <label>URL Server</label>
    <input type="text" name="server_url" value="{{.Cfg.ServerURL}}">
    <label>Token Agent</label>
    <input type="text" name="agent_token" value="{{.Cfg.AgentToken}}">
    <label>Interval Polling (detik)</label>
    <input type="number" name="poll_interval_sec" value="{{.Cfg.PollIntervalSec}}">
    <label>Interval Heartbeat (detik)</label>
    <input type="number" name="heartbeat_sec" value="{{.Cfg.HeartbeatSec}}">
    <button type="submit">Simpan</button>
  </form>
</div>

</body>
</html>
`

type dashboardData struct {
	Saved         bool
	IsHealthy     bool
	StartedAt     string
	LastHeartbeat string
	LastPoll      string
	RecentJobs    []JobEvent
	Cfg           *Config
}

// startWebUI jalankan dashboard lokal di 127.0.0.1:<port>.
// Sengaja cuma bind ke localhost (bukan 0.0.0.0) supaya tidak bisa
// diakses dari luar PC itu — tidak perlu auth tambahan.
func (p *program) startWebUI() {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		snap := p.state.Snapshot()
		p.cfgMu.RLock()
		cfg := *p.cfg
		p.cfgMu.RUnlock()

		data := dashboardData{
			Saved:         r.URL.Query().Get("saved") == "1",
			IsHealthy:     !snap.LastHeartbeat.IsZero(),
			StartedAt:     snap.StartedAt.Format("2006-01-02 15:04:05"),
			LastHeartbeat: formatTimeOrDash(snap.LastHeartbeat),
			LastPoll:      formatTimeOrDash(snap.LastPoll),
			RecentJobs:    snap.RecentJobs,
			Cfg:           &cfg,
		}
		_ = tmpl.Execute(w, data)
	})

	mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		r.ParseForm()

		p.cfgMu.Lock()
		p.cfg.ServerURL = r.FormValue("server_url")
		p.cfg.AgentToken = r.FormValue("agent_token")
		if v, err := strconv.Atoi(r.FormValue("poll_interval_sec")); err == nil && v > 0 {
			p.cfg.PollIntervalSec = v
		}
		if v, err := strconv.Atoi(r.FormValue("heartbeat_sec")); err == nil && v > 0 {
			p.cfg.HeartbeatSec = v
		}
		cfgCopy := *p.cfg
		p.cfgMu.Unlock()

		// Simpan ke file + bikin ulang API client pakai setting baru,
		// jadi langsung aktif tanpa perlu restart service
		_ = cfgCopy.Save(p.configPath)
		p.apiMu.Lock()
		p.api = NewAPIClient(&cfgCopy)
		p.apiMu.Unlock()

		http.Redirect(w, r, "/?saved=1", http.StatusSeeOther)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", p.cfg.DashboardPort)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Warningf("dashboard web gagal start di %s: %v", addr, err)
		}
	}()
}

func formatTimeOrDash(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

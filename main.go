package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	cfg    *Config
	api    *APIClient
	stopCh chan struct{}
}

func (p *program) Start(s service.Service) error {
	p.stopCh = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(p.stopCh)
	return nil
}

func (p *program) run() {
	go p.heartbeatLoop()
	p.pollLoop()
}

// pollLoop cek job baru tiap N detik + jitter kecil biar tidak barengan
// dengan agent lain saat request ke server
func (p *program) pollLoop() {
	tempDir := os.TempDir()

	for {
		select {
		case <-p.stopCh:
			return
		default:
		}

		jobs, err := p.api.FetchJobs()
		if err != nil {
			logger.Warningf("gagal ambil job: %v", err)
		}

		for _, job := range jobs {
			p.processJob(job, tempDir)
		}

		jitter := time.Duration(rand.Intn(5)) * time.Second
		interval := time.Duration(p.cfg.PollIntervalSec)*time.Second + jitter
		time.Sleep(interval)
	}
}

func (p *program) processJob(job PrintJob, tempDir string) {
	_ = p.api.UpdateJobStatus(job.ID, "printing", "")

	filePath, err := downloadFile(job.FileURL, tempDir, job.ID)
	if err != nil {
		logger.Errorf("job %d: gagal download: %v", job.ID, err)
		_ = p.api.UpdateJobStatus(job.ID, "failed", err.Error())
		return
	}
	defer cleanupFile(filePath)

	if err := printPDF(filePath, job.PrinterName, job.Copies); err != nil {
		logger.Errorf("job %d: gagal print: %v", job.ID, err)
		_ = p.api.UpdateJobStatus(job.ID, "failed", err.Error())
		return
	}

	_ = p.api.UpdateJobStatus(job.ID, "done", "")
	logger.Infof("job %d selesai dicetak ke %s", job.ID, job.PrinterName)
}

// heartbeatLoop lapor "masih hidup" secara terpisah dari polling job,
// intervalnya lebih santai (default 60 detik)
func (p *program) heartbeatLoop() {
	ticker := time.NewTicker(time.Duration(p.cfg.HeartbeatSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			if err := p.api.Heartbeat(); err != nil {
				logger.Warningf("heartbeat gagal: %v", err)
			}
		}
	}
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("gagal baca config.json: %v", err)
	}

	svcConfig := &service.Config{
		Name:        "MarsaPrintAgent",
		DisplayName: "MARSA Print Agent",
		Description: "Agent cetak jarak jauh untuk MARSA Remote Print",
	}

	prg := &program{
		cfg: cfg,
		api: NewAPIClient(cfg),
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Jalankan sebagai: agent.exe install  → daftar sebagai service
	//                   agent.exe start    → jalankan service
	//                   agent.exe stop     → hentikan service
	//                   agent.exe uninstall → hapus service
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

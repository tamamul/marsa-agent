package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewAPIClient(cfg *Config) *APIClient {
	return &APIClient{
		baseURL: cfg.ServerURL,
		token:   cfg.AgentToken,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *APIClient) doRequest(method, path string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server balas status %d: %s", resp.StatusCode, string(respBody))
	}
	return respBody, nil
}

// PrintJob merepresentasikan satu job cetak dari server
type PrintJob struct {
	ID          int    `json:"id"`
	PrinterName string `json:"printer_name"` // nama printer persis sesuai driver Windows
	FileURL     string `json:"file_url"`     // link download file (idealnya sudah PDF)
	Copies      int    `json:"copies"`
}

// Heartbeat lapor "masih hidup" ke server
func (c *APIClient) Heartbeat() error {
	_, err := c.doRequest("POST", "/api/agent/heartbeat", nil)
	return err
}

// FetchJobs ambil job pending yang jadi tanggung jawab agent ini
func (c *APIClient) FetchJobs() ([]PrintJob, error) {
	respBody, err := c.doRequest("GET", "/api/agent/jobs", nil)
	if err != nil {
		return nil, err
	}
	var jobs []PrintJob
	if err := json.Unmarshal(respBody, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// UpdateJobStatus lapor status job: printing / done / failed
func (c *APIClient) UpdateJobStatus(jobID int, status string, errMsg string) error {
	payload, _ := json.Marshal(map[string]string{
		"status":        status,
		"error_message": errMsg,
	})
	_, err := c.doRequest("POST", fmt.Sprintf("/api/agent/jobs/%d/status", jobID), payload)
	return err
}

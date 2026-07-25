package n8n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebhookPayload struct {
	Event     string    `json:"event"`
	Token     string    `json:"token"`
	IP        string    `json:"ip"`
	Status    string    `json:"status"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

// DispatchWebhook sends an async JSON payload to the specified n8n Webhook URL.
// If the URL is empty or unreachable, it falls back gracefully with a simulated log.
func DispatchWebhook(webhookURL string, payload interface{}) (bool, string) {
	webhookURL = strings.TrimSpace(webhookURL)
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Sprintf("JSON marshal error: %v", err)
	}

	log.Printf("[n8n Webhook Dispatcher] Event triggered -> Payload: %s", string(jsonBytes))

	if webhookURL == "" {
		return true, "Simulated Webhook Dispatch (No webhook URL configured)"
	}

	client := http.Client{
		Timeout: 4 * time.Second,
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return false, fmt.Sprintf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "AgencyPulse-n8n-Dispatcher/0.6.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[n8n Webhook Warning] Delivery to %s failed (%v), fallback to simulated dispatch", webhookURL, err)
		return true, fmt.Sprintf("Webhook attempt dispatched (Fallback mode, endpoint error: %v)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, fmt.Sprintf("Webhook successfully delivered (HTTP %d)", resp.StatusCode)
	}

	return true, fmt.Sprintf("Webhook delivered with status (HTTP %d)", resp.StatusCode)
}

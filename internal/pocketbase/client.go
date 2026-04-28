package pocketbase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Collections esperadas en PocketBase:
// industries, solutions, case_studies, prospects, proposals, events.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewFromEnv() *Client {
	return &Client{
		baseURL: os.Getenv("POCKETBASE_URL"),
		token:   os.Getenv("POCKETBASE_ADMIN_TOKEN"),
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c.baseURL != ""
}

func (c *Client) CreateRecord(ctx context.Context, collection string, payload map[string]any) error {
	if !c.Enabled() {
		return nil
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/api/collections/%s/records", c.baseURL, collection), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("pocketbase create %s status %d", collection, resp.StatusCode)
	}
	return nil
}

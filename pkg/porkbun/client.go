package porkbun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.porkbun.com/api/json/v3"

type Client struct {
	APIKey       string
	SecretAPIKey string
	httpClient   *http.Client
}

func NewClient(apiKey, secretAPIKey string) *Client {
	return &Client{
		APIKey:       apiKey,
		SecretAPIKey: secretAPIKey,
		httpClient:   &http.Client{},
	}
}

// authBody is embedded in every request
type authBody struct {
	APIKey       string `json:"apikey"`
	SecretAPIKey string `json:"secretapikey"`
}

type DNSRecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"` // subdomain, empty for root
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
}

type createDNSRequest struct {
	authBody
	Name    string `json:"name,omitempty"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
}

type dnsResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	ID      json.Number `json:"id,omitempty"`
}

type listDNSResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Records []DNSRecord `json:"records,omitempty"`
}

type deleteRequest struct {
	authBody
}

func (c *Client) post(endpoint string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	resp, err := c.httpClient.Post(baseURL+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	return nil
}

// CreateRecord creates a DNS record for the given domain.
// name is the subdomain (empty string for root).
func (c *Client) CreateRecord(domain, name, recordType, content, ttl string) (string, error) {
	req := createDNSRequest{
		authBody: authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey},
		Name:     name,
		Type:     recordType,
		Content:  content,
		TTL:      ttl,
	}

	var result dnsResponse
	if err := c.post("/dns/create/"+domain, req, &result); err != nil {
		return "", err
	}

	if result.Status != "SUCCESS" {
		return "", fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.ID.String(), nil
}

// DeleteRecord deletes a DNS record by domain and record ID.
func (c *Client) DeleteRecord(domain, id string) error {
	req := deleteRequest{
		authBody: authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey},
	}

	var result dnsResponse
	if err := c.post("/dns/delete/"+domain+"/"+id, req, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("porkbun error: %s", result.Message)
	}

	return nil
}

// Ping verifies credentials and returns your public IP as seen by Porkbun.
func (c *Client) Ping() (string, error) {
	req := authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey}

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
		YourIP  string `json:"yourIp,omitempty"`
	}

	if err := c.post("/ping", req, &result); err != nil {
		return "", err
	}

	if result.Status != "SUCCESS" {
		return "", fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.YourIP, nil
}

// ListRecords returns all DNS records for the given domain.
func (c *Client) ListRecords(domain string) ([]DNSRecord, error) {
	req := authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey}

	var result listDNSResponse
	if err := c.post("/dns/retrieve/"+domain, req, &result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" {
		return nil, fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.Records, nil
}
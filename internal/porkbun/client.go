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

type editDNSRequest struct {
	authBody
	Name    string `json:"name,omitempty"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     string `json:"ttl"`
}

type Domain struct {
	Domain     string `json:"domain"`
	Status     string `json:"status"`
	TLD        string `json:"tld"`
	CreateDate string `json:"createDate"`
	ExpireDate string `json:"expireDate"`
	AutoRenew  any    `json:"autoRenew"`
}

type listDomainsResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message,omitempty"`
	Domains []Domain `json:"domains,omitempty"`
}

type TLDPricing struct {
	Registration string `json:"registration"`
	Renewal      string `json:"renewal"`
	Transfer     string `json:"transfer"`
}

type pricingResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message,omitempty"`
	Pricing map[string]TLDPricing `json:"pricing,omitempty"`
}

type SSLBundle struct {
	CertificateChain string `json:"certificatechain"`
	PrivateKey       string `json:"privatekey"`
	PublicKey        string `json:"publickey"`
}

type sslResponse struct {
	Status           string `json:"status"`
	Message          string `json:"message,omitempty"`
	CertificateChain string `json:"certificatechain,omitempty"`
	PrivateKey       string `json:"privatekey,omitempty"`
	PublicKey        string `json:"publickey,omitempty"`
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

// EditRecord edits a DNS record by domain and record ID.
func (c *Client) EditRecord(domain, id, name, recordType, content, ttl string) error {
	req := editDNSRequest{
		authBody: authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey},
		Name:     name,
		Type:     recordType,
		Content:  content,
		TTL:      ttl,
	}

	var result dnsResponse
	if err := c.post("/dns/edit/"+domain+"/"+id, req, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("porkbun error: %s", result.Message)
	}

	return nil
}

// EditRecordByType edits DNS records by domain, type, and optional subdomain.
func (c *Client) EditRecordByType(domain, recordType, subdomain, content, ttl string) error {
	endpoint := "/dns/editByNameType/" + domain + "/" + recordType
	if subdomain != "" {
		endpoint += "/" + subdomain
	}

	req := editDNSRequest{
		authBody: authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey},
		Content:  content,
		Type:     recordType,
		TTL:      ttl,
	}

	var result dnsResponse
	if err := c.post(endpoint, req, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("porkbun error: %s", result.Message)
	}

	return nil
}

// DeleteRecordByType deletes DNS records by domain, type, and optional subdomain.
func (c *Client) DeleteRecordByType(domain, recordType, subdomain string) error {
	endpoint := "/dns/deleteByNameType/" + domain + "/" + recordType
	if subdomain != "" {
		endpoint += "/" + subdomain
	}

	req := deleteRequest{
		authBody: authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey},
	}

	var result dnsResponse
	if err := c.post(endpoint, req, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("porkbun error: %s", result.Message)
	}

	return nil
}

// ListRecordsByType returns DNS records for a domain filtered by type and optional subdomain.
func (c *Client) ListRecordsByType(domain, recordType, subdomain string) ([]DNSRecord, error) {
	endpoint := "/dns/retrieveByNameType/" + domain + "/" + recordType
	if subdomain != "" {
		endpoint += "/" + subdomain
	}

	req := authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey}

	var result listDNSResponse
	if err := c.post(endpoint, req, &result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" {
		return nil, fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.Records, nil
}

// ListDomains returns all domains in the account.
func (c *Client) ListDomains() ([]Domain, error) {
	req := authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey}

	var result listDomainsResponse
	if err := c.post("/domain/listAll", req, &result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" {
		return nil, fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.Domains, nil
}

// GetPricing returns pricing for all TLDs. Returns a map of TLD to pricing info.
func (c *Client) GetPricing() (map[string]TLDPricing, error) {
	req := struct{}{}

	var result pricingResponse
	if err := c.post("/pricing/get", req, &result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" {
		return nil, fmt.Errorf("porkbun error: %s", result.Message)
	}

	return result.Pricing, nil
}

// RetrieveSSL retrieves the SSL certificate bundle for a domain.
func (c *Client) RetrieveSSL(domain string) (*SSLBundle, error) {
	req := authBody{APIKey: c.APIKey, SecretAPIKey: c.SecretAPIKey}

	var result sslResponse
	if err := c.post("/ssl/retrieve/"+domain, req, &result); err != nil {
		return nil, err
	}

	if result.Status != "SUCCESS" {
		return nil, fmt.Errorf("porkbun error: %s", result.Message)
	}

	return &SSLBundle{
		CertificateChain: result.CertificateChain,
		PrivateKey:       result.PrivateKey,
		PublicKey:        result.PublicKey,
	}, nil
}
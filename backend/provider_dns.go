package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ProviderDNSZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status,omitempty"`
}

type ProviderDNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
	Proxied *bool  `json:"proxied,omitempty"`
}

func (a *HTTPProviderAdapter) cloudflareRequest(ctx context.Context, method, path string, query url.Values, body io.Reader, out any) error {
	if !a.Enabled() { return fmt.Errorf("%s adapter is not configured", a.ID()) }
	u := strings.TrimRight(a.BaseURL, "/") + path
	if len(query) > 0 { u += "?" + query.Encode() }
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil { return err }
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Token)
	resp, err := a.Client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil { return err }
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return fmt.Errorf("cloudflare API: %s: %s", resp.Status, strings.TrimSpace(string(data))) }
	if out == nil || len(data) == 0 { return nil }
	if err := json.Unmarshal(data, out); err != nil { return fmt.Errorf("cloudflare API decode: %w", err) }
	return nil
}

type cloudflareEnvelope[T any] struct {
	Success bool `json:"success"`
	Result  T    `json:"result"`
}

func (a *HTTPProviderAdapter) CloudflareZones(ctx context.Context) ([]ProviderDNSZone, error) {
	if a.ProviderID != AdapterCloudflare { return nil, fmt.Errorf("provider is not cloudflare") }
	var env cloudflareEnvelope[[]ProviderDNSZone]
	q := url.Values{}
	q.Set("per_page", "100")
	if err := a.cloudflareRequest(ctx, http.MethodGet, "/zones", q, nil, &env); err != nil { return nil, err }
	if !env.Success { return nil, fmt.Errorf("cloudflare returned success=false") }
	return env.Result, nil
}

func (a *HTTPProviderAdapter) CloudflareZoneID(ctx context.Context, zone string) (string, error) {
	zone = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(zone)), ".")
	if zone == "" { return "", fmt.Errorf("zone is required") }
	zones, err := a.CloudflareZones(ctx)
	if err != nil { return "", err }
	for _, z := range zones { if strings.EqualFold(strings.TrimSuffix(z.Name, "."), zone) { return z.ID, nil } }
	return "", fmt.Errorf("cloudflare zone not found: %s", zone)
}

func (a *HTTPProviderAdapter) CloudflareRecords(ctx context.Context, zone string) ([]ProviderDNSRecord, error) {
	zoneID, err := a.CloudflareZoneID(ctx, zone)
	if err != nil { return nil, err }
	var env cloudflareEnvelope[[]ProviderDNSRecord]
	q := url.Values{}
	q.Set("per_page", "100")
	if err := a.cloudflareRequest(ctx, http.MethodGet, "/zones/"+url.PathEscape(zoneID)+"/dns_records", q, nil, &env); err != nil { return nil, err }
	if !env.Success { return nil, fmt.Errorf("cloudflare returned success=false") }
	return env.Result, nil
}

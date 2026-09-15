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

type CloudflareZoneSummary struct { ID string `json:"id"`; Name string `json:"name"`; Status string `json:"status,omitempty"` }
type CloudflareDNSRecord struct { ID string `json:"id,omitempty"`; ZoneID string `json:"zone_id,omitempty"`; ZoneName string `json:"zone_name,omitempty"`; Name string `json:"name"`; Type string `json:"type"`; Content string `json:"content,omitempty"`; TTL int `json:"ttl,omitempty"`; Proxied *bool `json:"proxied,omitempty"`; Comment string `json:"comment,omitempty"` }
type cloudflareAPIResponse struct { Success bool `json:"success"`; Result json.RawMessage `json:"result"` }

func (a *HTTPProviderAdapter) cloudflareRequest(ctx context.Context, method, path string, body []byte) (json.RawMessage,error) {
 if !a.Enabled(){return nil,fmt.Errorf("cloudflare adapter is not configured")}
 var reader io.Reader; if body!=nil{reader=strings.NewReader(string(body))}
 req,err:=http.NewRequestWithContext(ctx,method,strings.TrimRight(a.BaseURL,"/")+path,reader);if err!=nil{return nil,err}
 req.Header.Set("Accept","application/json");req.Header.Set("Authorization","Bearer "+a.Token);if body!=nil{req.Header.Set("Content-Type","application/json")}
 resp,err:=a.Client.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();raw,err:=io.ReadAll(io.LimitReader(resp.Body,4<<20));if err!=nil{return nil,err};if resp.StatusCode<200||resp.StatusCode>=300{return nil,fmt.Errorf("cloudflare API: %s",resp.Status)}
 var out cloudflareAPIResponse;if err:=json.Unmarshal(raw,&out);err!=nil{return nil,err};if !out.Success{return nil,fmt.Errorf("cloudflare API returned unsuccessful response")};return out.Result,nil
}
func (a *HTTPProviderAdapter) CloudflareZones(ctx context.Context)([]CloudflareZoneSummary,error){raw,err:=a.cloudflareRequest(ctx,http.MethodGet,"/zones?per_page=100",nil);if err!=nil{return nil,err};var v []CloudflareZoneSummary;if err:=json.Unmarshal(raw,&v);err!=nil{return nil,err};return v,nil}
func (a *HTTPProviderAdapter) CloudflareZoneID(ctx context.Context,zone string)(string,error){q:=url.Values{};q.Set("name",strings.TrimSuffix(strings.ToLower(strings.TrimSpace(zone)),"."));q.Set("status","active");q.Set("per_page","5");raw,err:=a.cloudflareRequest(ctx,http.MethodGet,"/zones?"+q.Encode(),nil);if err!=nil{return "",err};var v []CloudflareZoneSummary;if err:=json.Unmarshal(raw,&v);err!=nil{return "",err};for _,z:=range v{if strings.EqualFold(z.Name,strings.TrimSuffix(zone,".")){return z.ID,nil}};return "",fmt.Errorf("cloudflare zone not found: %s",zone)}
func (a *HTTPProviderAdapter) CloudflareRecords(ctx context.Context,zone string)([]CloudflareDNSRecord,error){id,err:=a.CloudflareZoneID(ctx,zone);if err!=nil{return nil,err};raw,err:=a.cloudflareRequest(ctx,http.MethodGet,"/zones/"+url.PathEscape(id)+"/dns_records?per_page=100",nil);if err!=nil{return nil,err};var v []CloudflareDNSRecord;if err:=json.Unmarshal(raw,&v);err!=nil{return nil,err};return v,nil}
func (a *HTTPProviderAdapter) CloudflareRecordCreate(ctx context.Context,zone string,record json.RawMessage)(json.RawMessage,error){id,err:=a.CloudflareZoneID(ctx,zone);if err!=nil{return nil,err};return a.cloudflareRequest(ctx,http.MethodPost,"/zones/"+url.PathEscape(id)+"/dns_records",record)}
func (a *HTTPProviderAdapter) CloudflareRecordUpdate(ctx context.Context,zone,recordID string,record json.RawMessage)(json.RawMessage,error){id,err:=a.CloudflareZoneID(ctx,zone);if err!=nil{return nil,err};if strings.TrimSpace(recordID)==""{return nil,fmt.Errorf("record_id is required for update")};return a.cloudflareRequest(ctx,http.MethodPut,"/zones/"+url.PathEscape(id)+"/dns_records/"+url.PathEscape(recordID),record)}
func (a *HTTPProviderAdapter) CloudflareRecordDelete(ctx context.Context,zone,recordID string)(json.RawMessage,error){id,err:=a.CloudflareZoneID(ctx,zone);if err!=nil{return nil,err};if strings.TrimSpace(recordID)==""{return nil,fmt.Errorf("record_id is required for delete")};return a.cloudflareRequest(ctx,http.MethodDelete,"/zones/"+url.PathEscape(id)+"/dns_records/"+url.PathEscape(recordID),nil)}

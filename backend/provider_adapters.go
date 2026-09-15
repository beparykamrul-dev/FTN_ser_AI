package backend

import (
    "context"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)

// HTTPProviderAdapter is the common live adapter for external providers.
// Each provider keeps its credentials outside the FTN source tree.
type HTTPProviderAdapter struct {
    ProviderID AdapterKind
    BaseURL string
    HealthPath string
    Token string
    Client *http.Client
}

func NewHTTPProviderAdapter(id AdapterKind, baseEnv, pathEnv, tokenEnv, defaultURL, defaultPath string) *HTTPProviderAdapter {
    base:=strings.TrimRight(os.Getenv(baseEnv),"/"); if base==""{base=defaultURL}
    path:=os.Getenv(pathEnv); if path==""{path=defaultPath}
    return &HTTPProviderAdapter{ProviderID:id,BaseURL:base,HealthPath:path,Token:os.Getenv(tokenEnv),Client:&http.Client{Timeout:10*time.Second}}
}
func (a *HTTPProviderAdapter) ID() string{return string(a.ProviderID)}
func (a *HTTPProviderAdapter) Kind() AdapterKind{return a.ProviderID}
func (a *HTTPProviderAdapter) Enabled() bool{return strings.TrimSpace(a.Token)!=""}
func (a *HTTPProviderAdapter) Health(ctx context.Context) error{
    if !a.Enabled(){return errors.New("adapter is not configured")}
    req,err:=http.NewRequestWithContext(ctx,http.MethodGet,a.BaseURL+a.HealthPath,nil);if err!=nil{return err}
    req.Header.Set("Accept","application/json");req.Header.Set("Authorization","Bearer "+a.Token)
    resp,err:=a.Client.Do(req);if err!=nil{return err};defer resp.Body.Close()
    _,_=io.Copy(io.Discard,io.LimitReader(resp.Body,1<<20))
    if resp.StatusCode<200||resp.StatusCode>=300{return fmt.Errorf("%s adapter: %s",a.ID(),resp.Status)}
    return nil
}

func NewCloudflareAdapterFromEnv()*HTTPProviderAdapter{return NewHTTPProviderAdapter(AdapterCloudflare,"FTN_CLOUDFLARE_API_URL","FTN_CLOUDFLARE_HEALTH_PATH","FTN_CLOUDFLARE_TOKEN","https://api.cloudflare.com","/user/tokens/verify")}
func NewGoogleAdapterFromEnv()*HTTPProviderAdapter{return NewHTTPProviderAdapter(AdapterGoogle,"FTN_GOOGLE_API_URL","FTN_GOOGLE_HEALTH_PATH","FTN_GOOGLE_TOKEN","https://dns.googleapis.com","/dns/v1/projects")}
func NewRenderAdapterFromEnv()*HTTPProviderAdapter{return NewHTTPProviderAdapter(AdapterRender,"FTN_RENDER_API_URL","FTN_RENDER_HEALTH_PATH","FTN_RENDER_API_KEY","https://api.render.com","/v1/services")}

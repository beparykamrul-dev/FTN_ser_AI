package backend

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)

type GitHubAdapter struct { BaseURL string; Token string; Owner string; Repo string; Branch string; Client *http.Client }
func NewGitHubAdapterFromEnv() *GitHubAdapter { base:=strings.TrimRight(os.Getenv("FTN_GITHUB_API_URL"),"/"); if base==""{base="https://api.github.com"}; branch:=os.Getenv("FTN_GITHUB_BRANCH"); if branch==""{branch="feat/ftndns-authority-anycast-v2"}; return &GitHubAdapter{BaseURL:base,Token:os.Getenv("FTN_GITHUB_TOKEN"),Owner:os.Getenv("FTN_GITHUB_OWNER"),Repo:os.Getenv("FTN_GITHUB_REPO"),Branch:branch,Client:&http.Client{Timeout:15*time.Second}} }
func (g *GitHubAdapter) ID() string { return "github" }
func (g *GitHubAdapter) Kind() AdapterKind { return AdapterGitHub }
func (g *GitHubAdapter) Enabled() bool { return g.Token!="" && g.Owner!="" && g.Repo!="" }
func (g *GitHubAdapter) Health(ctx context.Context) error { _,_,err:=g.request(ctx,http.MethodGet,"/repos/"+g.Owner+"/"+g.Repo,nil); return err }
func (g *GitHubAdapter) request(ctx context.Context,method,path string,body io.Reader)([]byte,int,error){if !g.Enabled(){return nil,0,errors.New("github adapter is not configured")}; req,err:=http.NewRequestWithContext(ctx,method,g.BaseURL+path,body);if err!=nil{return nil,0,err};req.Header.Set("Accept","application/vnd.github+json");req.Header.Set("X-GitHub-Api-Version","2022-11-28");req.Header.Set("Authorization","Bearer "+g.Token);if body!=nil{req.Header.Set("Content-Type","application/json")};resp,err:=g.Client.Do(req);if err!=nil{return nil,0,err};defer resp.Body.Close();data,err:=io.ReadAll(io.LimitReader(resp.Body,4<<20));if err!=nil{return nil,resp.StatusCode,err};if resp.StatusCode<200||resp.StatusCode>=300{return data,resp.StatusCode,fmt.Errorf("github api: %s",resp.Status)};return data,resp.StatusCode,nil}
func (g *GitHubAdapter) Repository(ctx context.Context)(map[string]any,error){data,_,err:=g.request(ctx,http.MethodGet,"/repos/"+g.Owner+"/"+g.Repo,nil);if err!=nil{return nil,err};var out map[string]any;if err:=json.Unmarshal(data,&out);err!=nil{return nil,err};return out,nil}
func (g *GitHubAdapter) DispatchWorkflow(ctx context.Context,workflow string,inputs map[string]string)error{if strings.TrimSpace(workflow)==""||strings.ContainsAny(workflow,"/\\")||strings.Contains(workflow,".."){return errors.New("invalid workflow")};body,err:=json.Marshal(map[string]any{"ref":g.Branch,"inputs":inputs});if err!=nil{return err};_,_,err=g.request(ctx,http.MethodPost,"/repos/"+g.Owner+"/"+g.Repo+"/actions/workflows/"+workflow+"/dispatches",strings.NewReader(string(body)));return err}
func (g *GitHubAdapter) VerifyWebhook(payload []byte,signature string)bool{secret:=os.Getenv("FTN_GITHUB_WEBHOOK_SECRET");if secret==""||!strings.HasPrefix(signature,"sha256="){return false};mac:=hmac.New(sha256.New,[]byte(secret));_,_=mac.Write(payload);expected:="sha256="+hex.EncodeToString(mac.Sum(nil));return hmac.Equal([]byte(expected),[]byte(signature))}

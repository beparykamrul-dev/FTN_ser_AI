package backend

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
)

func (s *APIServer) githubStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet { w.Header().Set("Allow", http.MethodGet); writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error":"method not allowed"}); return }
    writeJSON(w, http.StatusOK, map[string]any{
        "enabled": s.github.Enabled(),
        "adapter": "github",
        "mode": "live-adapter",
        "repository": map[string]string{"owner":s.github.Owner,"repo":s.github.Repo,"branch":s.github.Branch},
        "features": []string{"repository-read","workflow-dispatch","webhook-events"},
    })
}

func (s *APIServer) githubRepository(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet { w.Header().Set("Allow", http.MethodGet); writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error":"method not allowed"}); return }
    repo, err := s.github.Repository(r.Context())
    if err != nil { writeJSON(w, http.StatusBadGateway, map[string]any{"ok":false,"error":err.Error()}); return }
    writeJSON(w, http.StatusOK, map[string]any{"ok":true,"repository":repo})
}

func (s *APIServer) githubWebhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { w.Header().Set("Allow", http.MethodPost); writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error":"method not allowed"}); return }
    body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 2<<20))
    if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid webhook body"}); return }
    if !s.github.VerifyWebhook(body, r.Header.Get("X-Hub-Signature-256")) { writeJSON(w, http.StatusUnauthorized, map[string]string{"error":"invalid webhook signature"}); return }
    var event map[string]any
    if err := json.Unmarshal(body, &event); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid JSON"}); return }
    writeJSON(w, http.StatusAccepted, map[string]any{"ok":true,"event":r.Header.Get("X-GitHub-Event"),"action":event["action"],"queued":true})
}

func (s *APIServer) githubWorkflow(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { w.Header().Set("Allow", http.MethodPost); writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error":"method not allowed"}); return }
    guard := strings.TrimSpace(r.Header.Get("X-FTN-Control-Token"))
    expected := strings.TrimSpace(getEnv("FTN_CONTROL_SECRET"))
    if expected == "" || guard == "" || !secureEqual(guard, expected) { writeJSON(w, http.StatusUnauthorized, map[string]string{"error":"invalid control credentials"}); return }
    var in struct { Workflow string `json:"workflow"`; Inputs map[string]string `json:"inputs"` }
    dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 128<<10)); dec.DisallowUnknownFields()
    if err := dec.Decode(&in); err != nil || strings.TrimSpace(in.Workflow) == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"workflow is required"}); return }
    if strings.Contains(in.Workflow, "/") || strings.Contains(in.Workflow, "\\") || strings.Contains(in.Workflow, "..") { writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid workflow"}); return }
    if err := s.github.DispatchWorkflow(r.Context(), in.Workflow, in.Inputs); err != nil { writeJSON(w, http.StatusBadGateway, map[string]any{"ok":false,"error":err.Error()}); return }
    writeJSON(w, http.StatusAccepted, map[string]any{"ok":true,"workflow":in.Workflow,"branch":s.github.Branch,"status":"dispatched"})
}

func getEnv(key string) string { return strings.TrimSpace(envLookup(key)) }

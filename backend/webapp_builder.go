package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UIComponent struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Props       map[string]interface{} `json:"props"`
	Children    []UIComponent `json:"children,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Page struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Components  []UIComponent `json:"components"`
	Meta        map[string]string `json:"meta"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WebApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Pages       []Page `json:"pages"`
	Theme       map[string]string `json:"theme"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var webApps = make(map[string]*WebApp)
var components = make(map[string]UIComponent)

var componentLibrary = map[string]map[string]interface{}{
	"button": {
		"name": "Button",
		"icon": "🔘",
		"props": map[string]string{"text": "string", "color": "string", "size": "string"},
	},
	"input": {
		"name": "Input Field",
		"icon": "📝",
		"props": map[string]string{"placeholder": "string", "type": "string"},
	},
	"card": {
		"name": "Card",
		"icon": "📋",
		"props": map[string]string{"title": "string", "description": "string"},
	},
	"header": {
		"name": "Header",
		"icon": "📌",
		"props": map[string]string{"title": "string", "subtitle": "string"},
	},
	"footer": {
		"name": "Footer",
		"icon": "🔚",
		"props": map[string]string{"content": "string"},
	},
	"navbar": {
		"name": "Navigation Bar",
		"icon": "📍",
		"props": map[string]string{"items": "array"},
	},
	"grid": {
		"name": "Grid",
		"icon": "⚙️",
		"props": map[string]string{"columns": "number", "gap": "string"},
	},
	"modal": {
		"name": "Modal",
		"icon": "📦",
		"props": map[string]string{"title": "string", "content": "string"},
	},
	"dropdown": {
		"name": "Dropdown",
		"icon": "▼",
		"props": map[string]string{"items": "array"},
	},
	"badge": {
		"name": "Badge",
		"icon": "🏷️",
		"props": map[string]string{"text": "string", "color": "string"},
	},
	"table": {
		"name": "Table",
		"icon": "📊",
		"props": map[string]string{"columns": "array", "rows": "array"},
	},
	"chart": {
		"name": "Chart",
		"icon": "📈",
		"props": map[string]string{"type": "string", "data": "array"},
	},
	"form": {
		"name": "Form",
		"icon": "📋",
		"props": map[string]string{"fields": "array", "submit": "string"},
	},
	"list": {
		"name": "List",
		"icon": "📝",
		"props": map[string]string{"items": "array"},
	},
	"accordion": {
		"name": "Accordion",
		"icon": "📂",
		"props": map[string]string{"items": "array"},
	},
	"tabs": {
		"name": "Tabs",
		"icon": "📑",
		"props": map[string]string{"tabs": "array"},
	},
	"alert": {
		"name": "Alert",
		"icon": "⚠️",
		"props": map[string]string{"message": "string", "type": "string"},
	},
	"spinner": {
		"name": "Spinner",
		"icon": "⏳",
		"props": map[string]string{"size": "string"},
	},
	"slider": {
		"name": "Slider",
		"icon": "🎚️",
		"props": map[string]string{"min": "number", "max": "number"},
	},
	"toggle": {
		"name": "Toggle",
		"icon": "🔘",
		"props": map[string]string{"checked": "boolean"},
	},
	"checkbox": {
		"name": "Checkbox",
		"icon": "☑️",
		"props": map[string]string{"label": "string"},
	},
	"radio": {
		"name": "Radio",
		"icon": "◯",
		"props": map[string]string{"options": "array"},
	},
	"progress": {
		"name": "Progress Bar",
		"icon": "📊",
		"props": map[string]string{"value": "number", "max": "number"},
	},
	"avatar": {
		"name": "Avatar",
		"icon": "👤",
		"props": map[string]string{"src": "string", "alt": "string"},
	},
	"image": {
		"name": "Image",
		"icon": "🖼️",
		"props": map[string]string{"src": "string", "alt": "string"},
	},
	"video": {
		"name": "Video",
		"icon": "🎥",
		"props": map[string]string{"src": "string"},
	},
	"text": {
		"name": "Text",
		"icon": "📄",
		"props": map[string]string{"content": "string", "size": "string"},
	},
	"divider": {
		"name": "Divider",
		"icon": "───",
		"props": map[string]string{},
	},
	"spacer": {
		"name": "Spacer",
		"icon": "⬌",
		"props": map[string]string{"height": "string"},
	},
	"container": {
		"name": "Container",
		"icon": "📦",
		"props": map[string]string{"width": "string", "padding": "string"},
	},
}

// Create Web App
func createWebApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webapp WebApp
	err := json.NewDecoder(r.Body).Decode(&webapp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	webapp.ID = fmt.Sprintf("app_%d", time.Now().Unix())
	webapp.CreatedAt = time.Now()
	webapp.UpdatedAt = time.Now()
	if webapp.Pages == nil {
		webapp.Pages = []Page{}
	}

	webApps[webapp.ID] = &webapp

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(webapp)
}

// Get Web App
func getWebApp(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("id")
	if appID == "" {
		http.Error(w, "App ID required", http.StatusBadRequest)
		return
	}

	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// Get all Web Apps
func getAllWebApps(w http.ResponseWriter, r *http.Request) {
	var apps []*WebApp
	for _, app := range webApps {
		apps = append(apps, app)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"apps": apps,
	})
}

// Add page to Web App
func addPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appID := r.URL.Query().Get("app_id")
	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	var page Page
	err := json.NewDecoder(r.Body).Decode(&page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	page.ID = fmt.Sprintf("page_%d", time.Now().Unix())
	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	app.Pages = append(app.Pages, page)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(page)
}

// Add component to page
func addComponent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appID := r.URL.Query().Get("app_id")
	pageID := r.URL.Query().Get("page_id")

	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	var pageIndex int = -1
	for i, p := range app.Pages {
		if p.ID == pageID {
			pageIndex = i
			break
		}
	}

	if pageIndex == -1 {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	var component UIComponent
	err := json.NewDecoder(r.Body).Decode(&component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component.ID = fmt.Sprintf("comp_%d", time.Now().Unix())
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()

	app.Pages[pageIndex].Components = append(app.Pages[pageIndex].Components, component)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(component)
}

// Get component library
func getComponentLibrary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"components": componentLibrary,
		"total": len(componentLibrary),
	})
}

// Export as HTML
func exportHTML(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("id")
	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	html := generateHTML(app)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.html", app.ID))
	fmt.Fprint(w, html)
}

// Export as JSON
func exportJSON(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("id")
	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", app.ID))
	json.NewEncoder(w).Encode(app)
}

// Update Web App
func updateWebApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appID := r.URL.Query().Get("id")
	app, exists := webApps[appID]
	if !exists {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	var updates WebApp
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Name = updates.Name
	app.Description = updates.Description
	app.Theme = updates.Theme
	app.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// Delete Web App
func deleteWebApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appID := r.URL.Query().Get("id")
	delete(webApps, appID)

	w.WriteHeader(http.StatusNoContent)
}

func generateHTML(app *WebApp) string {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #333; }
		.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
		.card { background: white; padding: 20px; margin: 15px 0; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
		h1 { color: white; text-align: center; margin: 30px 0; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
		p { line-height: 1.6; color: #666; }
		.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
		.button { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 12px 24px; border: none; border-radius: 8px; cursor: pointer; font-weight: 600; }
		.button:hover { transform: translateY(-2px); box-shadow: 0 6px 15px rgba(102, 126, 234, 0.4); }
	</style>
</head>
<body>
	<div class="container">
		<h1>🎨 %s</h1>
		<p style="text-align: center; color: white; margin-bottom: 40px;">%s</p>
`, app.Name, app.Name, app.Description)

	for _, page := range app.Pages {
		html += fmt.Sprintf(`
		<div class="card">
			<h2>📄 %s</h2>
			<div class="grid">
`, page.Name)

		for _, comp := range page.Components {
			html += fmt.Sprintf(`
			<div class="card">
				<h3>%s</h3>
				<p>Component: %s</p>
			</div>
`, comp.Name, comp.Type)
		}

		html += `
			</div>
		</div>
`
	}

	html += `
	</div>
</body>
</html>
`
	return html
}

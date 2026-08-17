package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type TimeZoneClock struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TimeZone    string `json:"timezone"`
	CurrentTime string `json:"current_time"`
}

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var db *sql.DB
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan interface{}, 100)
var mu sync.Mutex
var startTime time.Time

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	startTime = time.Now()
	var err error
	connStr := "user=postgres password=yourpassword dbname=clock_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Database connection error:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ Database ping error:", err)
	}

	createTables()
	log.Println("✅ Database connected!")
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS timezones (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		timezone VARCHAR(50) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Table creation error:", err)
	}
}

func main() {
	mux := http.NewServeMux()

	// WebSocket
	mux.HandleFunc("/ws", handleWebSocket)

	// REST APIs
	mux.HandleFunc("/api/timezones", getTimeZoneClocks)
	mux.HandleFunc("/api/timezone/add", addTimeZone)
	mux.HandleFunc("/api/timezone/delete", deleteTimeZone)
	mux.HandleFunc("/api/health", getServerHealth)
	mux.HandleFunc("/api/dashboard", getDashboardData)

	// Maintenance Endpoints
	mux.HandleFunc("/api/maintenance/tasks", getAllMaintenanceTasks)
	mux.HandleFunc("/api/maintenance/vacuum", runDatabaseVacuum)
	mux.HandleFunc("/api/maintenance/reindex", runIndexRebuild)
	mux.HandleFunc("/api/maintenance/logs", runLogCleanup)
	mux.HandleFunc("/api/maintenance/cache", runCacheClear)
	mux.HandleFunc("/api/maintenance/pool", runConnectionPoolReset)
	mux.HandleFunc("/api/maintenance/backup", runBackupDatabase)
	mux.HandleFunc("/api/maintenance/cleanup", runSystemCleanup)
	mux.HandleFunc("/api/maintenance/memory", runMemoryOptimization)
	mux.HandleFunc("/api/maintenance/security", runSecurityAudit)
	mux.HandleFunc("/api/maintenance/tuning", runPerformanceTuning)
	mux.HandleFunc("/api/maintenance/backups", getBackups)

	// Health check
	mux.HandleFunc("/api/health", healthCheck)

	// CORS
	c := cors.Default()
	handler := c.Handler(mux)

	// Broadcast goroutine
	go broadcastUpdates()
	go updateCurrentTimes()

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func broadcastUpdates() {
	for msg := range broadcast {
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("❌ Error writing to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func updateCurrentTimes() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		clocks := fetchAllClocks()
		broadcast <- WebSocketMessage{
			Type: "update",
			Data: clocks,
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	log.Printf("✅ New WebSocket connection. Total clients: %d", len(clients))

	clocks := fetchAllClocks()
	conn.WriteJSON(WebSocketMessage{
		Type: "init",
		Data: clocks,
	})

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			log.Printf("❌ Client disconnected. Total clients: %d", len(clients))
			break
		}

		handleMessage(msg)
	}
}

func handleMessage(msg WebSocketMessage) {
	switch msg.Type {
	case "add":
		data := msg.Data.(map[string]interface{})
		addTimeZoneWS(data)
	case "delete":
		data := msg.Data.(map[string]interface{})
		deleteTimeZoneWS(data)
	}
}

func addTimeZoneWS(data map[string]interface{}) {
	name := data["name"].(string)
	timezone := data["timezone"].(string)

	var id int
	err := db.QueryRow(
		"INSERT INTO timezones (name, timezone) VALUES ($1, $2) RETURNING id",
		name, timezone,
	).Scan(&id)

	if err != nil {
		broadcast <- WebSocketMessage{
			Type: "error",
			Data: err.Error(),
		}
		return
	}

	clock := TimeZoneClock{
		ID:          id,
		Name:        name,
		TimeZone:    timezone,
		CurrentTime: getCurrentTime(timezone),
	}

	broadcast <- WebSocketMessage{
		Type: "add",
		Data: clock,
	}
}

func deleteTimeZoneWS(data map[string]interface{}) {
	id := int(data["id"].(float64))

	_, err := db.Exec("DELETE FROM timezones WHERE id = $1", id)
	if err != nil {
		broadcast <- WebSocketMessage{
			Type: "error",
			Data: err.Error(),
		}
		return
	}

	broadcast <- WebSocketMessage{
		Type: "delete",
		Data: map[string]int{"id": id},
	}
}

func fetchAllClocks() []TimeZoneClock {
	rows, err := db.Query("SELECT id, name, timezone FROM timezones ORDER BY created_at DESC")
	if err != nil {
		log.Printf("❌ Error fetching clocks: %v", err)
		return []TimeZoneClock{}
	}
	defer rows.Close()

	var clocks []TimeZoneClock

	for rows.Next() {
		var clock TimeZoneClock
		err := rows.Scan(&clock.ID, &clock.Name, &clock.TimeZone)
		if err != nil {
			log.Println(err)
			continue
		}

		clock.CurrentTime = getCurrentTime(clock.TimeZone)
		clocks = append(clocks, clock)
	}

	return clocks
}

func getTimeZoneClocks(w http.ResponseWriter, r *http.Request) {
	clocks := fetchAllClocks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clocks": clocks,
	})
}

func addTimeZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Name     string `json:"name"`
		TimeZone string `json:"timezone"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow(
		"INSERT INTO timezones (name, timezone) VALUES ($1, $2) RETURNING id",
		data.Name, data.TimeZone,
	).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clock := TimeZoneClock{
		ID:          id,
		Name:        data.Name,
		TimeZone:    data.TimeZone,
		CurrentTime: getCurrentTime(data.TimeZone),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(clock)
}

func deleteTimeZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM timezones WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func getCurrentTime(timezone string) string {
	var tm time.Time
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "Invalid"
	}
	tm = time.Now().In(loc)
	return tm.Format("15:04:05")
}

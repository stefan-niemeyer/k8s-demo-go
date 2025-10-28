package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	port       = getEnv("PORT", "3000")
	appVersion = getEnv("APP_VERSION", "v1")
	appPicture = getEnv("APP_PICTURE", "v1.jpg")
	unstable   = getEnv("UNSTABLE", "")
	host       = getEnv("HOSTNAME", "localhost")
	pathBase   = "/"
	pathState  = "/state"
	pathReady  = "/ready"
	pathHealth = "/health"
	pathCrash  = "/crash"

	startTime     = time.Now()
	totalRequests int
	readyState    = true
	healthState   = true
)

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func createResponse(endpoint, method string) map[string]interface{} {
	log.Printf("Host: %s | Path: %s | Version: %s | Total Requests: %d | Ready: %v | Health: %v | App Uptime: %.2f seconds | Log Time: %s\n",
		host, endpoint, appVersion, totalRequests, readyState, healthState, time.Since(startTime).Seconds(), time.Now())
	return map[string]interface{}{
		"host":          host,
		"method":        method,
		"path":          endpoint,
		"version":       appVersion,
		"totalRequests": totalRequests,
		"readyState":    readyState,
		"healthState":   healthState,
	}
}

func main() {
	http.HandleFunc(pathBase, imageHandler)
	http.HandleFunc(pathState, stateHandler)
	http.HandleFunc(pathReady, readyHandler)
	http.HandleFunc(pathHealth, healthHandler)
	http.HandleFunc(pathCrash, crashHandler)

	if unstable != "" {
		if s, err := strconv.Atoi(unstable); err == nil && s > 0 {
			go func() {
				time.Sleep(time.Duration(s) * time.Second)
				log.Printf("server: UNSTABLE=%ss", unstable)
				healthState = false
			}()
		}
	}

	log.Printf("server: App listening on port %s!", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Returns the image file
func imageHandler(w http.ResponseWriter, r *http.Request) {
	if !healthState {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !readyState {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	totalRequests++
	createResponse(pathBase, "GET")
	imgPath := filepath.Join(".", "images", appPicture)
	http.ServeFile(w, r, imgPath)
}

// Returns the state of the app
func stateHandler(w http.ResponseWriter, r *http.Request) {
	if !healthState {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !readyState {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	totalRequests++
	resp := createResponse(pathState, "GET")
	json.NewEncoder(w).Encode(resp)
}

// Readiness probe GET/PUT
func readyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		statusCode := http.StatusOK
		if !readyState {
			statusCode = http.StatusServiceUnavailable
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]bool{"state": readyState})
		return
	}
	if r.Method == http.MethodPut {
		var data map[string]bool
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data["state"] == false && data["state"] == true { //ist false oder true
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		readyState = data["state"]
		resp := createResponse(pathReady, "PUT")
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Liveness probe GET/PUT
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		statusCode := http.StatusOK
		if !healthState {
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]bool{"state": healthState})
		return
	}
	if r.Method == http.MethodPut {
		var data map[string]bool
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data["state"] == false && data["state"] == true {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		healthState = data["state"]
		resp := createResponse(pathHealth, "PUT")
		json.NewEncoder(w).Encode(resp)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Crash endpoint
func crashHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Host: %s | Path: %s | Total Requests: %d | Health State: %v | Ready State: %v | App Uptime: %.2f seconds | Log Time: %s\n", host, pathCrash, totalRequests, healthState, readyState, time.Since(startTime).Seconds(), time.Now())
	http.Error(w, "Unexpected Error Occurred", http.StatusInternalServerError)
	os.Exit(1)
}

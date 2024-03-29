package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var counter int
var rwMutex = &sync.RWMutex{}
var lastModified time.Time
var version string

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	rwMutex.Lock()
	counter++
	lastModified = time.Now()
	rwMutex.Unlock()
	fmt.Fprintf(w, strconv.Itoa(counter))
}

func readCounter(w http.ResponseWriter, r *http.Request) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
	fmt.Fprintf(w, strconv.Itoa(counter))
}

func getVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Version: %s", version)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Page Not Found")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method Not Allowed")
		return
	}
	// Handle POST request here
	fmt.Fprintf(w, "Handling POST request")
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC3339)
	fmt.Fprintf(w, "Current time: %s", currentTime)
}

func handleQueryString(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	age := query.Get("age")
	fmt.Fprintf(w, "Name: %s, Age: %s", name, age)
}

func handleJSONResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"name": "John", "age": 30, "city": "New York"}`)
}

func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/file.txt")
}

func handleSecureEndpoint(w http.ResponseWriter, r *http.Request) {
	// Basic Authentication
	username := "user"
	password := "password"

	// Get the Authorization header from the request
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	// Check if the Authorization header is valid
	auth := strings.SplitN(authHeader, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}

	// Decode the base64-encoded credentials
	decoded, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}

	// Convert the decoded credentials to string
	credentials := string(decoded)

	// Compare the credentials with the expected username and password
	if credentials != fmt.Sprintf("%s:%s", username, password) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	// If authentication is successful, proceed with handling the request
	fmt.Fprintf(w, "This is a secure endpoint.")
}

func handleSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers here
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", echoString)
	r.HandleFunc("/increment", incrementCounter)
	r.HandleFunc("/counter", readCounter)
	r.HandleFunc("/version", getVersion)
	r.HandleFunc("/post", handlePost).Methods(http.MethodPost)
	r.HandleFunc("/time", handleTime)
	r.HandleFunc("/query", handleQueryString)
	r.HandleFunc("/json", handleJSONResponse)
	r.HandleFunc("/static", handleStaticFiles)
	r.HandleFunc("/secure", handleSecureEndpoint).Use(handleSecureHeaders) // Applying secure headers middleware
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)

	log.Fatal(http.ListenAndServe(":8081", r))
}

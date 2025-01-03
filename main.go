package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

// Map to store shortened URLs
var urls = make(map[string]string)

// Handles the form for submitting URLs
func handleForm(w http.ResponseWriter, r *http.Request) {
	// Log the request path for debugging
	fmt.Println("Handling request for:", r.URL.Path)

	// Handle POST requests to redirect to shortening endpoint
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	// Render the HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>URL Shortener</title>
	</head>
	<body>
		<h2>URL Shortener</h2>
		<form method="post" action="/shorten">
			<input type="url" name="url" placeholder="Enter a URL" required>
			<input type="submit" value="Shorten">
		</form>
	</body>
	</html>
	`)
}

// Handles URL shortening
func handleshorten(w http.ResponseWriter, r *http.Request) {
	// Log the request path for debugging
	fmt.Println("Handling request for:", r.URL.Path)

	// Validate POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the original URL
	originalUrl := r.FormValue("url")
	if originalUrl == "" {
		http.Error(w, "Missing url", http.StatusBadRequest)
		return
	}

	// Generate a short key and store the mapping
	shortkey := generateShortKey()
	urls[shortkey] = originalUrl

	// Dynamically create the shortened URL based on deployment host
	shortenedUrl := fmt.Sprintf("https://%s/short/%s", r.Host, shortkey)

	// Return the shortened URL
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>URL Shortener</title>
	</head>
	<body>
		<h2>URL Shortener</h2>
		<p>Original URL: %s</p>
		<p>Shortened URL: <a href="%s">%s</a></p>
	</body>
	</html>
	`, originalUrl, shortenedUrl, shortenedUrl)
}

// Handles redirection based on short key
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Log the request path for debugging
	fmt.Println("Handling redirect for:", r.URL.Path)

	// Get the short key from URL path
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")

	// Validate key
	if shortKey == "" {
		http.Error(w, "Key missing", http.StatusBadRequest)
		return
	}

	// Check if the key exists
	originalUrl, found := urls[shortKey]
	if !found {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Redirect to the original URL
	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
}

// Generates a random 6-character short key
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// Middleware to remove trailing slashes
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Normalize URL path by removing trailing slash
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, "/"), http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Main function to start the server
func main() {
	// Print environment and debug logs
	fmt.Println("Starting server...")

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleForm)
	mux.HandleFunc("/shorten", handleshorten)
	mux.HandleFunc("/short/", handleRedirect)

	// Wrap with trailing slash handler
	handler := removeTrailingSlash(mux)

	// Port binding
	port := os.Getenv("PORT")
	if port == "" {
		port = "4020" // Default port for local testing
	}

	// Start the server
	fmt.Println("Server running on port:", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

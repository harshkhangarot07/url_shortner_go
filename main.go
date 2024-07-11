package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

func handleshorten(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalUrl := r.FormValue("url")

	if originalUrl == "" {
		http.Error(w, "Missing url", http.StatusBadRequest)
		return
	}

	shortkey := generateShortKey()
	urls[shortkey] = originalUrl

	shortenedUrl := fmt.Sprintf("http://localhost:4020/short/%s", shortkey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<p>Original URL: `, originalUrl, `</p>
			<p>Shortened URL: <a href="`, shortenedUrl, `">`, shortenedUrl, `</a></p>
		</body>
		</html>
	`)

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]

	if shortKey == "" {
		http.Error(w, "key missing", http.StatusBadRequest)
		return
	}

	originalUrl, found := urls[shortKey]
	if !found {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)

}
func handleForm(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)

		return
	}

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

var urls = make(map[string]string)

func main() {

	// rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleshorten)
	http.HandleFunc("/short/", handleRedirect)

	fmt.Println("url shortener running on :4020")
	http.ListenAndServe(":4020", nil)
}

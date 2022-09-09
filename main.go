package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Configuration

// The address to bind the snippetbox srver to listen to.
// Generally you don't need to specify a host in the address unless your
// compter has multiple network interfaces and you just want to listen
// on one of them.
// NOTE: Should follow format "host:port".
const BindAddress = ":4000"

// HTTP Handler functions

// index is a catch-all handler routed at "/"
func index(w http.ResponseWriter, r *http.Request) {
	// Check if the request URL path exactly matches "/". If it doesn't,
	// use the http.NotFound() function to send a 404 response to the
	// client.
	// Importantly, we then return from the handler, otherwise we would
	// also write the hello message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id param from the query string and
	// try to convert it to an integer. If it fails, send a 404.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	// Interpolate the wanted snippet id into the response
	fmt.Fprintf(w, "Display snippet #%d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != http.MethodPost {
		// Add an "Allow: POST" header to the response header map.
		// This must be called before w.WriteHeader() or w.Write().
		w.Header().Set("Allow", http.MethodPost)
		// Send a Method Not Allowed HTTP error.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet ..."))
}

// Main application

func main() {
	// Initialize a new servemux (aka router) to store our URL paths.
	mux := http.NewServeMux()

	// Go's servemux supports two different types of URL patterns: fixed
	// paths and subtree paths. Fixed paths don't end in /, subtree paths
	// do.
	// Subtree path patterns are matched (and handled) whenever the start
	// of a request URL matches the subtree path. Imagine a wildcard at
	// the end of them: / = /**, /static/ = /static/**.
	// NOTES:
	// - Longer URL patterns take precedence over shorter ones.
	// - Request URL paths are automatically sanitized and the user
	//   redirected (if needed).
	// - If a subtree path is registered and a request to it without a
	//   trailing slash is received, user will be sent a 301 Permanent
	//   Redirect to the subtree path with the slash added.
	// - URL patterns can contain hostnames like this:
	//   mux.HandleFunc("foo.vaino.lol/", fooHandler)
	//   mux.HandleFunc("bar.vaino.lol/", barHandler)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/new", snippetCreate)

	// Start a new web server with the given network address to listen
	// on and the servemux we just created. If http.ListenAndServe()
	// returns an error, we use the log.Fatal() function to log the
	// error message and exit. Note that any error returned by
	// http.ListenAndServe() is always non-nil.
	log.Println("Snippetbox server listening on " + BindAddress)
	err := http.ListenAndServe(BindAddress, mux)
	log.Fatal(err)
}

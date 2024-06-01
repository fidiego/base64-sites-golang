package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type ApiRequest struct {
	Content string `json:"content"`
}

type ApiErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ApiResponse struct {
	Base64String          string  `json:"base64_string"`
	Base64StringLength    int     `json:"base64_string_length"`
	ContentLength         int     `json:"content_length"`
	ExecutionTime         float64 `json:"execution_time"`
	MinifiedContentLength int     `json:"minified_content_length"`
	Success               bool    `json:"success"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // start timer

	// TODO: check if request is a post

	// parse JSON POST request
	decoder := json.NewDecoder(r.Body)
	var p ApiRequest
	err := decoder.Decode(&p)

	// Handle JSON Error
	if err != nil {
		log.Println("Bad Request - returning error message")
		// set headers
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		// compose error response payload
		errorResponse := ApiErrorResponse{
			Success: false,
			Message: "Expecting a JSON payload with a single attribute \"content\" .",
		}
		log.Println(errorResponse)
		// write error response to buffer
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Happy Path 😊
	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	//  minify HTML
	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minified, err := minifier.String("text/html", p.Content)
	// Handle minification error
	if err != nil {
		log.Println("Error during HTML Minification")
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		errorResponse := ApiErrorResponse{
			Success: false,
			Message: "An unexpected error occurred while minifying the content. Make sure you're passing valid HTML.",
		}
		// write error response to buffer
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// encode
	encoded := base64.StdEncoding.EncodeToString([]byte(p.Content))
	base64_string := fmt.Sprintf("data:text/html;base64,%s", encoded)
	// calculate lengths
	content_length := len(p.Content)
	base64_string_length := len(base64_string)
	minified_content_length := len(minified)
	// get elapsed time
	elapsed := time.Now().Sub(start)

	// compose response object
	jsonResponse := ApiResponse{
		Success:               true,
		Base64StringLength:    base64_string_length,
		Base64String:          base64_string,
		ExecutionTime:         elapsed.Seconds(),
		ContentLength:         content_length,
		MinifiedContentLength: minified_content_length,
	}
	// write response
	json.NewEncoder(w).Encode(jsonResponse)
}

type RenderContext struct {
	Content template.URL
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	contents, present := query["content"]
	if !present || len(contents) == 0 {
		// return error response
	}
	t, err := template.ParseFiles("templates/render.go.html")
	if err != nil {
		log.Println(err)
	}
	content := contents[0]
	context := RenderContext{Content: template.URL(content)}
	t.Execute(w, context)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.go.html", "templates/index.js")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, nil)
}

func main() {
	fmt.Println("Starting Base64 Site Service")

	// parse port from env var - use 8080 as a fallback if not set
	var portStr = os.Getenv("PORT")
	var port int32
	if portStr == "" {
		log.Println("PORT env not provided. Falling back to default 8080.")
		port = 8080
	} else {
		parsed, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatal("Error parsing PORT env var")
			os.Exit(1)
		} else {
			port = int32(parsed)
		}

	}

	fmt.Println(" -> loaded env vars")
	fmt.Printf("     - PORT=%d\n", port)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        nil,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// register routes
	fmt.Println(" -> prepared http.Server")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/render", renderHandler)

	fmt.Println(" -> registered routes")
	fmt.Printf(" -> starting on port:%d\n", port)

	log.Fatal(server.ListenAndServe())
}

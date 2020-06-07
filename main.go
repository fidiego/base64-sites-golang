package main

import "encoding/base64"
import "encoding/json"
import "fmt"
import "html/template"
import "log"
import "net/http"
import "os"
import "strings"
import "time"

type ApiRequest struct {
	Content string `json:content`
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
			Message: "Expecting a JSON payload with a single attribute \"html\" .",
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
	// TODO: minify HTML

	// encode
	encoded := base64.StdEncoding.EncodeToString([]byte(p.Content))
	base64_string := strings.Join([]string{"data:text/html;base64,", encoded}, "")
	// calculate lengths
	content_length := len(p.Content)
	base64_string_length := len(base64_string)
	// get elapsed time
	elapsed := time.Now().Sub(start)

	// compose response object
	jsonResponse := ApiResponse{
		Success:            true,
		Base64StringLength: base64_string_length,
		Base64String:       base64_string,
		ExecutionTime:      elapsed.Seconds(),
		ContentLength:      content_length,
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

	var port = os.Getenv("PORT") // TODO: validate, perhaps
	fmt.Println(" -> loaded env vars")
	fmt.Printf("     - PORT=%s\n", port)

	server := &http.Server{
		Addr:           strings.Join([]string{":", port}, ""), // TODO: better concat?
		Handler:        nil,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println(" -> prepared http.Server")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/render", renderHandler)
	fmt.Println(" -> registered routes")
	fmt.Printf(" -> starting on port:%s\n", port)
	log.Fatal(server.ListenAndServe())
}

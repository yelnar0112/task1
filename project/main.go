package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	respID       int
	requestData  sync.Map
	responseData sync.Map
)

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type Response struct {
	ID      int               `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int               `json:"length"`
}

func main() {
	http.HandleFunc("/proxy", Proxy)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// CORS headers setup
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			h.ServeHTTP(w, r)
		})
	}

	log.Fatal(http.ListenAndServe(":"+port, corsHandler(http.DefaultServeMux)))
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Decoding error req: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Get(req.URL)
	if err != nil {
		http.Error(w, "Error sending req: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	status := resp.StatusCode
	header := make(map[string]string)
	for key, values := range resp.Header {
		header[key] = values[0]
	}

	respID++
	response := Response{
		ID:      respID,
		Status:  status,
		Headers: header,
		Length:  len(body),
	}

	proxyResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	requestData.Store(respID, req)
	responseData.Store(respID, proxyResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(proxyResponse)
}

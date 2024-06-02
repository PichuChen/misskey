package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

var misskeyNodeBackend = "http://localhost:3000"
var backend2ListenPort = "3001"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirdect to port 3000
		requestURL := r.URL
		requestMethod := r.Method
		requestHeaders := r.Header
		slog.Info("Request received", "method", requestMethod, "url", requestURL, "headers", requestHeaders)
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading request body", "err", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest(requestMethod, misskeyNodeBackend+requestURL.String(), bytes.NewReader(requestBody))
		if err != nil {
			slog.Error("Error creating request", "err", err)
			return
		}

		req.Header = requestHeaders
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Error sending request", "err", err)
			return
		}

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading response body", "err", err)
			return
		}

		for key, value := range resp.Header {
			w.Header().Set(key, value[0])
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)

		slog.Info("Request sent")

	})
	http.ListenAndServe(":"+backend2ListenPort, mux)

}

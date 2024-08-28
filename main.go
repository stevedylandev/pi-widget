package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/events", handleSSE)

	server := &http.Server{Addr: ":4321"}

	go func() {
		fmt.Println("Server is running on http://localhost:4321")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}

//go:embed index.html
var indexHTML string

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	log.Println("SSE connection established")

	for {
		combinedStats, err := getStats()
		if err != nil {
			log.Printf("Error getting IPFS Stats: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		data, err := json.Marshal(combinedStats)
		if err != nil {
			log.Printf("Error marshaling IPFS stats: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", data)
		if err != nil {
			log.Printf("Error writing to response: %v", err)
			return
		}
		w.(http.Flusher).Flush()

		time.Sleep(1 * time.Second)
	}
}

package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"

    "github.com/gorilla/mux"
    "streamline.go/connectors"
)

type RequestPayload struct {
    Provider   string          `json:"provider"`
    ConfigPath string          `json:"configPath,omitempty"`
    Config     json.RawMessage `json:"config,omitempty"`
    Messages   []string        `json:"messages"`
}

func processMessages(w http.ResponseWriter, r *http.Request) {
    var payload RequestPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var configData []byte
    if payload.ConfigPath != "" {
        data, err := os.ReadFile(payload.ConfigPath)
        if err != nil {
            http.Error(w, "Failed to read config file", http.StatusBadRequest)
            return
        }
        configData = data
    } else if payload.Config != nil {
        configData = payload.Config
    }

    if len(configData) == 0 {
        http.Error(w, "Configuration required", http.StatusBadRequest)
        return
    }

    connector, err := connectors.NewConnector(payload.Provider, configData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Use a WaitGroup to ensure the goroutine completes
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer wg.Done()

        if err := connector.Connect(); err != nil {
            log.Printf("Failed to connect: %v", err)
            return
        }

        for _, msg := range payload.Messages {
            if err := connector.Send(msg); err != nil {
                log.Printf("Failed to send message: %v", err)
                continue
            }
            log.Printf("Sent message: %s", msg)
        }

        if err := connector.Close(); err != nil {
            log.Printf("Failed to close connector: %v", err)
        }
    }()

    // Wait briefly to catch immediate errors, but don't block the HTTP response
    go func() {
        wg.Wait()
        log.Println("Message processing completed")
    }()

    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{"status": "processing"})
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/process", processMessages).Methods("POST")

    log.Println("Server starting on :8081")
    if err := http.ListenAndServe(":8081", router); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
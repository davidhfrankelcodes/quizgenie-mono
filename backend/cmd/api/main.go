// backend/cmd/api/main.go
package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  // A simple ping endpoint
  http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message":"pong"}`))
  })

  addr := ":8080"
  fmt.Println("Starting API server on", addr)
  if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

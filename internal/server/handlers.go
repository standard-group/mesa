package server

import "net/http"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Mesa!"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

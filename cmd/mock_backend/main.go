package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/token/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access": "mock-access-token",
			"refresh": "mock-refresh-token",
		})
	})

	mux.HandleFunc("/api/inyeccion-salesforce/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Just mock a success for now, can be extended for E2E tests
		json.NewEncoder(w).Encode(map[string]interface{}{
			"StatusCode": 200,
			"Mensaje": "Success",
			"NumeroSolicitud": 12345,
		})
	})

	log.Println("Starting Mock Backend on 0.0.0.0:8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", mux))
}

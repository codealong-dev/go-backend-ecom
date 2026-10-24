package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Error(w http.ResponseWriter, status int, message string) {
	type errorEnvelope struct {
		Error string `json:"error"`
	}

	writeJSON(w, status, errorEnvelope{Error: message})
}

func JSON(w http.ResponseWriter, status int, data any) {
	type dataEnvelope struct {
		Data any `json:"data"`
	}

	writeJSON(w, status, dataEnvelope{Data: data})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println("failed to write JSON response:", err)
	}
}

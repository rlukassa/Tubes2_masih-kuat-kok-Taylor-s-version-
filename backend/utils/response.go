// File ini adalah utilitas untuk response standar API di backend Little Alchemy 2.
// Berisi struct dan fungsi untuk mengirim response JSON yang konsisten.

package utils

import (
    "encoding/json" // Untuk encode JSON
    "net/http"      // Untuk kebutuhan HTTP response
)

// Response adalah struct standar untuk response API
type Response struct {
    Status  string      `json:"status"`             // Status response (success/error)
    Message string      `json:"message"`            // Pesan response
    Data    interface{} `json:"data,omitempty"`     // Data (opsional)
}

// SendResponse mengirim response JSON ke client
func SendResponse(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json") // Set header tipe konten
    w.WriteHeader(statusCode)                          // Set status code

    response := Response{
        Status:  status,   // Status response
        Message: message,  // Pesan response
        Data:    data,     // Data (jika ada)
    }

    json.NewEncoder(w).Encode(response) // Encode dan kirim response JSON
}
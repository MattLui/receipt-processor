package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"receipt-processor/models"
	"receipt-processor/points"
)

var receipts = make(map[string]int)

func ProcessReceiptHandler(w http.ResponseWriter, r *http.Request) {
    var receipt models.Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    points := points.CalculatePoints(receipt)
    id := uuid.New().String()
    receipts[id] = points

    response := map[string]string{"id": id}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func GetPointsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    points, found := receipts[id]
    if !found {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    response := map[string]int{"points": points}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

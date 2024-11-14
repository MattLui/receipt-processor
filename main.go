package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
    Retailer      string   `json:"retailer"`
    PurchaseDate  string   `json:"purchaseDate"`
    PurchaseTime  string   `json:"purchaseTime"`
    Items         []Item   `json:"items"`
    Total         string   `json:"total"`
}

type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

var receipts = make(map[string]int)

func main() {
		router := mux.NewRouter()
    router.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
    router.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
	}

	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || len(receipt.Items) == 0 || receipt.Total == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseFloat(receipt.Total, 64); err != nil {
		http.Error(w, "Invalid format for Total", http.StatusBadRequest)
		return
	}

	for _, item := range receipt.Items {
		if _, err := strconv.ParseFloat(item.Price, 64); err != nil {
			http.Error(w, "Invalid format for Item Price", http.StatusBadRequest)
			return
		}
	}

	points := calculatePoints(receipt)
	id := uuid.New().String()
	receipts[id] = points

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
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

func calculatePoints(receipt Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(re.FindAllString(receipt.Retailer, -1))

	// 50 points if the total is a round dollar amount with no cents.
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
			points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if math.Mod(total, 0.25) == 0 {
			points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Points for items with description lengths that are multiples of 3.
	for _, item := range receipt.Items {
			if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
					price, _ := strconv.ParseFloat(item.Price, 64)
					points += int(math.Ceil(price * 0.2))
			}
	}

	// 6 points if the day in the purchase date is odd.
	day, _ := strconv.Atoi(receipt.PurchaseDate[8:10])
	if day%2 != 0 {
			points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTime.Hour() == 14 || (purchaseTime.Hour() == 15 && purchaseTime.Minute() < 60) {
			points += 10
	}

	return points
}
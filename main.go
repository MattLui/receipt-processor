package main

import (
	"fmt"
	"log"
	"net/http"

	"receipt-processor/handlers"

	"github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/receipts/process", handlers.ProcessReceiptHandler).Methods("POST")
    router.HandleFunc("/receipts/{id}/points", handlers.GetPointsHandler).Methods("GET")

    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

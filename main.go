package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type Receipt struct {
	ID           string   `json:"id"`
	Retailer     string   `json:"retailer" binding:"required"`
	PurchaseDate string   `json:"purchaseDate" binding:"required"`
	PurchaseTime string   `json:"purchaseTime" binding:"required"`
	Total        string   `json:"total" binding:"required"`
	Items        []Item   `json:"items" binding:"required,min=1"`
	Points       int      `json:"points"`
}

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

var receipts map[string]Receipt


func main() {
	receipts = make(map[string]Receipt)

	router := gin.Default()
	router.GET("/receipts/list", getReceipts)
	router.POST("/receipts/process", gin.WrapF(processReceipt))
	router.GET("/receipts/:receipt_id/points", getReceiptPoints)

	port := "8080" 
	log.Fatal(router.Run(":" + port))
}

func getReceipts(c *gin.Context) {
	c.JSON(http.StatusOK, receipts)
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		errorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := binding.Validator.ValidateStruct(receipt); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	receipt.ID = uuid.New().String()

	points, err := calculatePoints(receipt)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	receipt.Points = points
	receipts[receipt.ID] = receipt

	response := struct {
		ID string `json:"id"`
	}{
		ID: receipt.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getReceiptPoints(c *gin.Context) {
	receiptID := c.Param("receipt_id")
	receipt, found := receipts[receiptID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"ERROR":   "Receipt Not Found",
			"MESSAGE": fmt.Sprintf("Receipt ID %s does not exist.", receiptID),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"points": receipt.Points,
	})
}
	
func calculatePoints(receipt Receipt) (int, error) {
	pointsTotal := 0

	for _, char := range receipt.Retailer {
		if isAlphanumeric(char) {
			pointsTotal++
		}
	}

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid total value: %v", err)
	}
	if total == math.Trunc(total) {
		pointsTotal += 50
	}

	if math.Mod(total, 0.25) == 0 {
		pointsTotal += 25
	}

	pointsTotal += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, fmt.Errorf("Invalid price value: %v", err)
			}
			pointsTotal += int(math.Ceil(price * 0.2))
		}
	}

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		return 0, fmt.Errorf("Invalid purchase date: %v", err)
	}
	if purchaseDate.Day()%2 != 0 {
		pointsTotal += 6
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return 0, fmt.Errorf("Invalid purchase time: %v", err)
	}
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
		pointsTotal += 10
	}

	return pointsTotal, nil
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

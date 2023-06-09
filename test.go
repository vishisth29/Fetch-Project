package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetReceipts(t *testing.T) {
	router := gin.Default()
	router.GET("/receipts/list", getReceipts)

	req, _ := http.NewRequest("GET", "/receipts/list", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{}\n", rec.Body.String())
}

func TestProcessReceipt(t *testing.T) {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)

	payload := `{
		"retailer": "Example Retailer",
		"purchaseDate": "2023-06-08",
		"purchaseTime": "15:30",
		"total": "25.00",
		"items": [
			{
				"shortDescription": "Item 1",
				"price": "10.00"
			},
			{
				"shortDescription": "Item 2",
				"price": "15.00"
			}
		]
	}`

	req, _ := http.NewRequest("POST", "/receipts/process", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "{\"id\":\"<generated-receipt-id>\"}\n", rec.Body.String())
}

func TestGetReceiptPoints(t *testing.T) {
	router := gin.Default()
	router.GET("/receipts/:receipt_id/points", getReceiptPoints)

	req, _ := http.NewRequest("GET", "/receipts/12345/points", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "{\"error\":\"Receipt Not Found\",\"message\":\"Receipt ID 12345 does not exist.\"}\n", rec.Body.String())
}

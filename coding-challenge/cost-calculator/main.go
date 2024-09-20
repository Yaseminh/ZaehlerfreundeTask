package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"io/ioutil"
	"log"
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
)

// Struct for meter readings
type MeterReading struct {
	Timestamp int64   `json:"timestamp"` // Timestamp in milliseconds
	Value     float64 `json:"value"`     // Consumption in kWh
}

// Struct for Awattar market prices
type MarketPrice struct {
	StartTimestamp int64   `json:"start_timestamp"`
	EndTimestamp   int64   `json:"end_timestamp"`
	Price          float64 `json:"marketprice"` // Price in Euro/MWh
}

// Struct to handle Awattar API response
type MarketDataResponse struct {
	Data []MarketPrice `json:"data"`
}

// Struct to handle incoming meter readings in the request
type MeterReadingRequest struct {
	Readings []MeterReading `json:"readings"`
}

// Fetch market prices from Awattar API
func fetchMarketPrices(startTimestamp, endTimestamp int64) ([]MarketPrice, error) {
	url := fmt.Sprintf("https://api.awattar.de/v1/marketdata?start=%d&end=%d", startTimestamp, endTimestamp)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var marketData MarketDataResponse
	err = json.Unmarshal(body, &marketData)
	if err != nil {
		return nil, err
	}

	return marketData.Data, nil
}

// Calculate energy cost using meter readings and market prices
func calculateEnergyCost(meterReadings []MeterReading, marketPrices []MarketPrice) (float64, error) {
	totalCost := 0.0

	for i := 0; i < len(meterReadings)-1; i++ {
		readingA := meterReadings[i]
		readingB := meterReadings[i+1]

		// Calculate time duration (ms) and energy consumption (kWh) between readings
		energyConsumption := readingB.Value - readingA.Value

		// Loop through market prices to find the corresponding price
		for _, price := range marketPrices {
			// Ensure reading timestamp falls within the market price interval
			if readingA.Timestamp >= price.StartTimestamp && readingA.Timestamp < price.EndTimestamp {
				// Convert price from Euro/MWh to Euro/kWh (divide by 1000)
				pricePerKwh := price.Price / 1000

				// Calculate cost for energy consumption during this interval
				totalCost += energyConsumption * pricePerKwh
				break
			}
		}
	}

	return totalCost, nil
}

// Main function to handle the POST request
func main() {
	r := gin.Default()
	// Enable CORS for all origins
	r.Use(cors.Default())

	// REST endpoint to handle energy cost calculation
	r.POST("/energy_cost", func(c *gin.Context) {
		var requestData MeterReadingRequest

		// Parse JSON request body
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure there are at least two meter readings to calculate the cost
		if len(requestData.Readings) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least two meter readings are required"})
			return
		}

		// Get start and end timestamps for the meter readings
		startTimestamp := requestData.Readings[0].Timestamp
		endTimestamp := requestData.Readings[len(requestData.Readings)-1].Timestamp

		// Fetch market prices from Awattar API
		marketPrices, err := fetchMarketPrices(startTimestamp, endTimestamp)
		if err != nil {
			log.Printf("Error fetching market prices: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch market prices"})
			return
		}

		// Calculate total energy cost based on meter readings and market prices
		totalCost, err := calculateEnergyCost(requestData.Readings, marketPrices)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate energy cost"})
			return
		}

		// Return total energy cost in the response
		c.JSON(http.StatusOK, gin.H{"total_cost": totalCost})
	})

	// Start the server
	r.Run(":8080")
}

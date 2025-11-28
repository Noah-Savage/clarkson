package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)


// Hammond export format
type HammondVehicle struct {
	Name       string `json:"name"`
	Make       string `json:"make"`
	Model      string `json:"model"`
	Year       int    `json:"year"`
	Odometer   float64 `json:"odometer"`
}

type HammondFuelEntry struct {
	Date      string  `json:"date"`
	Odometer  float64 `json:"odometer"`
	Gallons   float64 `json:"gallons"`
	Price     float64 `json:"cost_per_unit"`
	TotalCost float64 `json:"total_cost"`
}

type HammondExport struct {
	Vehicles []HammondVehicle `json:"vehicles"`
	Fuel     []HammondFuelEntry `json:"fuel_entries"`
}

func (app *Application) importHammondDatabase(c *gin.Context) {
	userID := c.GetUint("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	body, err := io.ReadAll(src)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	var export HammondExport
	if err := json.Unmarshal(body, &export); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Hammond export format"})
		return
	}

	imported := gin.H{
		"vehicles": 0,
		"fuel":     0,
		"errors":   []string{},
	}

	// Import vehicles
	for _, v := range export.Vehicles {
		vehicle := Vehicle{
			UserID:      userID,
			Make:        v.Make,
			Model:       v.Model,
			Year:        v.Year,
			Odometer:    v.Odometer,
			MileageUnit: "mi", // Hammond default
			FuelType:    "Petrol",
		}
		if err := app.db.Create(&vehicle).Error; err != nil {
			imported["errors"] = append(imported["errors"].([]string), fmt.Sprintf("Failed to import vehicle %s: %v", v.Name, err))
			continue
		}

		// Import fuel entries
		for _, fuel := range export.Fuel {
			date, _ := time.Parse("2006-01-02", fuel.Date)
			entry := FuelEntry{
				VehicleID: vehicle.ID,
				Date:      date,
				Odometer:  fuel.Odometer,
				Gallons:   fuel.Gallons,
				Price:     fuel.TotalCost,
			}
			if err := app.db.Create(&entry).Error; err != nil {
				imported["errors"] = append(imported["errors"].([]string), fmt.Sprintf("Failed to import fuel entry: %v", err))
				continue
			}
			imported["fuel"] = imported["fuel"].(int) + 1
		}
		imported["vehicles"] = imported["vehicles"].(int) + 1
	}

	c.JSON(200, imported)
}

func (app *Application) importFuellyCSV(c *gin.Context) {
	userID := c.GetUint("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// TODO: Implement CSV parsing for Fuelly format
	// Fuelly CSV format: Date, Odometer, Fuel Type, Gallons, Price, Cost, Car, Notes

	imported := gin.H{
		"vehicles": 0,
		"fuel":     0,
		"errors":   []string{},
	}

	c.JSON(200, imported)
}

func (app *Application) importClarksonBackup(c *gin.Context) {
	userID := c.GetUint("userID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	body, err := io.ReadAll(src)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid Clarkson backup format"})
		return
	}

	imported := gin.H{
		"vehicles":   0,
		"fuel":       0,
		"expenses":   0,
		"reminders":  0,
		"errors":     []string{},
	}

	// Parse and import vehicles
	if vehicles, ok := data["vehicles"].([]interface{}); ok {
		for _, v := range vehicles {
			if vmap, ok := v.(map[string]interface{}); ok {
				vehicle := Vehicle{
					UserID:      userID,
					Make:        vmap["make"].(string),
					Model:       vmap["model"].(string),
					Year:        int(vmap["year"].(float64)),
					Odometer:    vmap["odometer"].(float64),
					MileageUnit: vmap["mileage_unit"].(string),
					FuelType:    vmap["fuel_type"].(string),
				}
				if err := app.db.Create(&vehicle).Error; err != nil {
					imported["errors"] = append(imported["errors"].([]string), fmt.Sprintf("Failed to import vehicle: %v", err))
					continue
				}
				imported["vehicles"] = imported["vehicles"].(int) + 1
			}
		}
	}

	c.JSON(200, imported)
}

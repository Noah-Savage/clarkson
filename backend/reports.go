package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)


type FuelTrendPoint struct {
	Month     string  `json:"month"`
	Cost      float64 `json:"cost"`
	Gallons   float64 `json:"gallons"`
	Distance  float64 `json:"distance"`
	MPG       float64 `json:"mpg"`
}

type ExpenseTrendPoint struct {
	Month     string             `json:"month"`
	Total     float64            `json:"total"`
	Categories map[string]float64 `json:"categories"`
}

type VehicleReportData struct {
	Vehicle         Vehicle                   `json:"vehicle"`
	FuelEntries     []FuelEntry               `json:"fuel_entries"`
	Expenses        []Expense                 `json:"expenses"`
	FuelTrend       []FuelTrendPoint          `json:"fuel_trend"`
	ExpenseTrend    []ExpenseTrendPoint       `json:"expense_trend"`
	TotalCost       float64                   `json:"total_cost"`
	TotalDistance   float64                   `json:"total_distance"`
	AverageMPG      float64                   `json:"average_mpg"`
	FuelCosts       float64                   `json:"fuel_costs"`
	MaintenanceCost float64                   `json:"maintenance_cost"`
	OtherCosts      float64                   `json:"other_costs"`
}

func (app *Application) generateDetailedReport(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicle Vehicle
	if err := app.db.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Vehicle not found"})
		return
	}

	var fuelEntries []FuelEntry
	app.db.Where("vehicle_id = ?", vehicleID).Order("date ASC").Find(&fuelEntries)

	var expenses []Expense
	app.db.Where("vehicle_id = ?", vehicleID).Order("date ASC").Find(&expenses)

	report := VehicleReportData{
		Vehicle:      vehicle,
		FuelEntries:  fuelEntries,
		Expenses:     expenses,
	}

	// Calculate fuel statistics
	if len(fuelEntries) > 1 {
		for i := 1; i < len(fuelEntries); i++ {
			distance := fuelEntries[i].Odometer - fuelEntries[i-1].Odometer
			report.TotalDistance += distance
		}

		totalGallons := 0.0
		for _, f := range fuelEntries {
			report.FuelCosts += f.Price
			totalGallons += f.Gallons
		}

		if totalGallons > 0 {
			report.AverageMPG = report.TotalDistance / totalGallons
		}

		// Build fuel trend
		monthlyFuel := make(map[string]gin.H)
		for i := 0; i < len(fuelEntries)-1; i++ {
			month := fuelEntries[i].Date.Format("2006-01")
			distance := fuelEntries[i+1].Odometer - fuelEntries[i].Odometer

			if data, exists := monthlyFuel[month]; exists {
				data["cost"] = data["cost"].(float64) + fuelEntries[i].Price
				data["gallons"] = data["gallons"].(float64) + fuelEntries[i].Gallons
				data["distance"] = data["distance"].(float64) + distance
			} else {
				monthlyFuel[month] = gin.H{
					"month":    month,
					"cost":     fuelEntries[i].Price,
					"gallons":  fuelEntries[i].Gallons,
					"distance": distance,
				}
			}
		}

		for _, data := range monthlyFuel {
			gallons := data["gallons"].(float64)
			distance := data["distance"].(float64)
			mpg := 0.0
			if gallons > 0 {
				mpg = distance / gallons
			}

			trend := FuelTrendPoint{
				Month:    data["month"].(string),
				Cost:     data["cost"].(float64),
				Gallons:  gallons,
				Distance: distance,
				MPG:      mpg,
			}
			report.FuelTrend = append(report.FuelTrend, trend)
		}
	}

	// Calculate expense statistics
	maintenanceCost := 0.0
	otherCost := 0.0
	for _, e := range expenses {
		report.TotalCost += e.Amount
		if e.Category == "Maintenance" {
			maintenanceCost += e.Amount
		} else {
			otherCost += e.Amount
		}
	}
	report.MaintenanceCost = maintenanceCost
	report.OtherCosts = otherCost

	// Build expense trend
	monthlyExpense := make(map[string]ExpenseTrendPoint)
	for _, e := range expenses {
		month := e.Date.Format("2006-01")
		if data, exists := monthlyExpense[month]; exists {
			data.Total += e.Amount
			data.Categories[e.Category] += e.Amount
			monthlyExpense[month] = data
		} else {
			cats := make(map[string]float64)
			cats[e.Category] = e.Amount
			monthlyExpense[month] = ExpenseTrendPoint{
				Month:      month,
				Total:      e.Amount,
				Categories: cats,
			}
		}
	}

	for _, data := range monthlyExpense {
		report.ExpenseTrend = append(report.ExpenseTrend, data)
	}

	report.TotalCost += report.FuelCosts

	c.JSON(200, report)
}

func (app *Application) exportDetailedCSV(c *gin.Context) {
	userID := c.GetUint("userID")
	c.Header("Content-Disposition", "attachment; filename=clarkson-detailed-export.csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")

	csv := "CLARKSON DETAILED EXPORT\n"
	csv += fmt.Sprintf("Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	totalAllCosts := 0.0

	for _, v := range vehicles {
		csv += fmt.Sprintf("\nVEHICLE: %d %s %s\n", v.Year, v.Make, v.Model)
		csv += fmt.Sprintf("Current Odometer: %.1f %s\n", v.Odometer, v.MileageUnit)
		csv += "\n--- FUEL ENTRIES ---\n"
		csv += "Date,Odometer,Gallons,Price,Location\n"

		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Order("date ASC").Find(&fuelEntries)

		totalFuelCost := 0.0
		for _, f := range fuelEntries {
			csv += fmt.Sprintf("%s,%.1f,%.2f,%.2f,%s\n",
				f.Date.Format("2006-01-02"), f.Odometer, f.Gallons, f.Price, f.Location)
			totalFuelCost += f.Price
		}

		csv += "\n--- EXPENSES ---\n"
		csv += "Date,Category,Amount,Notes\n"

		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Order("date ASC").Find(&expenses)

		totalExpenseCost := 0.0
		for _, e := range expenses {
			csv += fmt.Sprintf("%s,%s,%.2f,\"%s\"\n",
				e.Date.Format("2006-01-02"), e.Category, e.Amount, e.Notes)
			totalExpenseCost += e.Amount
		}

		totalAllCosts += totalFuelCost + totalExpenseCost

		csv += fmt.Sprintf("\nVehicle Total: $%.2f\n", totalFuelCost+totalExpenseCost)
	}

	csv += fmt.Sprintf("\n\nGRAND TOTAL: $%.2f\n", totalAllCosts)

	c.Data(200, "text/csv; charset=utf-8", []byte(csv))
}

func (app *Application) exportJSON(c *gin.Context) {
	userID := c.GetUint("userID")
	c.Header("Content-Disposition", "attachment; filename=clarkson-backup.json")
	c.Header("Content-Type", "application/json")

	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	export := gin.H{
		"version":   "1.0",
		"exported":  time.Now().Format(time.RFC3339),
		"vehicles":  []gin.H{},
	}

	vehicleList := []gin.H{}

	for _, v := range vehicles {
		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Find(&fuelEntries)

		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Find(&expenses)

		var reminders []MaintenanceReminder
		app.db.Where("vehicle_id = ?", v.ID).Find(&reminders)

		vehicleData := gin.H{
			"vehicle":    v,
			"fuel":       fuelEntries,
			"expenses":   expenses,
			"reminders":  reminders,
		}

		vehicleList = append(vehicleList, vehicleData)
	}

	export["vehicles"] = vehicleList

	c.JSON(200, export)
}

func (app *Application) generateComparisonReport(c *gin.Context) {
	userID := c.GetUint("userID")

	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	type VehicleComparison struct {
		Vehicle       Vehicle `json:"vehicle"`
		TotalCost     float64 `json:"total_cost"`
		TotalMiles    float64 `json:"total_miles"`
		AverageMPG    float64 `json:"average_mpg"`
		CostPerMile   float64 `json:"cost_per_mile"`
		FuelCount     int     `json:"fuel_count"`
		ExpenseCount  int     `json:"expense_count"`
	}

	comparisons := []VehicleComparison{}
	totalCost := 0.0

	for _, v := range vehicles {
		comp := VehicleComparison{Vehicle: v}

		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Order("date ASC").Find(&fuelEntries)
		comp.FuelCount = len(fuelEntries)

		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Find(&expenses)
		comp.ExpenseCount = len(expenses)

		if len(fuelEntries) > 1 {
			firstOdo := fuelEntries[0].Odometer
			lastOdo := fuelEntries[len(fuelEntries)-1].Odometer
			comp.TotalMiles = lastOdo - firstOdo

			totalGallons := 0.0
			for _, f := range fuelEntries {
				comp.TotalCost += f.Price
				totalGallons += f.Gallons
			}

			if totalGallons > 0 {
				comp.AverageMPG = comp.TotalMiles / totalGallons
			}

			if comp.TotalMiles > 0 {
				comp.CostPerMile = comp.TotalCost / comp.TotalMiles
			}
		}

		for _, e := range expenses {
			comp.TotalCost += e.Amount
		}

		totalCost += comp.TotalCost
		comparisons = append(comparisons, comp)
	}

	c.JSON(200, gin.H{
		"vehicles":   comparisons,
		"total_cost": totalCost,
	})
}

func (app *Application) generateSearchResults(c *gin.Context) {
	userID := c.GetUint("userID")
	query := c.Query("q")

	if query == "" {
		c.JSON(400, gin.H{"error": "Query parameter required"})
		return
	}

	var fuelEntries []FuelEntry
	app.db.
		Joins("JOIN vehicles ON vehicles.id = fuel_entries.vehicle_id").
		Where("vehicles.user_id = ? AND (fuel_entries.location LIKE ? OR fuel_entries.notes LIKE ?)", userID, "%"+query+"%", "%"+query+"%").
		Find(&fuelEntries)

	var expenses []Expense
	app.db.
		Joins("JOIN vehicles ON vehicles.id = expenses.vehicle_id").
		Where("vehicles.user_id = ? AND (expenses.category LIKE ? OR expenses.notes LIKE ?)", userID, "%"+query+"%", "%"+query+"%").
		Find(&expenses)

	results := gin.H{
		"fuel":     fuelEntries,
		"expenses": expenses,
		"count":    len(fuelEntries) + len(expenses),
	}

	c.JSON(200, results)
}

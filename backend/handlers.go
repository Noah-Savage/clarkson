package main

import (
	"fmt"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


// Vehicle Handlers

func (app *Application) listVehiclesWithStats(c *gin.Context) {
	userID := c.GetUint("userID")
	var vehicles []Vehicle

	if err := app.db.
		Joins("LEFT JOIN vehicle_users ON vehicle_users.vehicle_id = vehicles.id").
		Where("vehicles.user_id = ? OR vehicle_users.user_id = ?", userID, userID).
		Distinct("vehicles.*").
		Find(&vehicles).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Enhance with stats
	type VehicleWithStats struct {
		Vehicle        Vehicle   `json:"vehicle"`
		TotalCost      float64   `json:"total_cost"`
		TotalMiles     float64   `json:"total_miles"`
		AverageMPG     float64   `json:"average_mpg"`
		FuelCount      int64     `json:"fuel_count"`
		ExpenseCount   int64     `json:"expense_count"`
		LastFuelDate   *time.Time `json:"last_fuel_date"`
		DueReminders   int       `json:"due_reminders"`
	}

	var results []VehicleWithStats

	for _, v := range vehicles {
		stats := VehicleWithStats{Vehicle: v}

		// Calculate fuel stats
		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Order("date ASC").Find(&fuelEntries)

		stats.FuelCount = int64(len(fuelEntries))
		totalGallons := 0.0
		var totalCost float64

		if len(fuelEntries) > 0 {
			firstOdometer := fuelEntries[0].Odometer
			lastOdometer := fuelEntries[len(fuelEntries)-1].Odometer
			stats.TotalMiles = lastOdometer - firstOdometer
			stats.LastFuelDate = &fuelEntries[len(fuelEntries)-1].Date

			for _, f := range fuelEntries {
				totalGallons += f.Gallons
				totalCost += f.Price
			}

			if totalGallons > 0 {
				stats.AverageMPG = stats.TotalMiles / totalGallons
			}
		}

		// Calculate expense stats
		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Find(&expenses)
		stats.ExpenseCount = int64(len(expenses))

		for _, e := range expenses {
			totalCost += e.Amount
		}

		stats.TotalCost = totalCost

		// Count due reminders
		var reminders []MaintenanceReminder
		app.db.Where("vehicle_id = ?", v.ID).Find(&reminders)

		for _, r := range reminders {
			status := "upcoming"
			if r.IntervalMiles > 0 {
				nextServiceMiles := r.LastServiceMiles + r.IntervalMiles
				if v.Odometer >= nextServiceMiles {
					status = "overdue"
				} else if v.Odometer >= nextServiceMiles-500 {
					status = "soon"
				}
			}

			if r.IntervalDays > 0 {
				nextServiceDate := r.LastServiceDate.AddDate(0, 0, r.IntervalDays)
				if time.Now().After(nextServiceDate) {
					status = "overdue"
				} else if time.Until(nextServiceDate).Hours()/24 < 7 {
					status = "soon"
				}
			}

			if status != "upcoming" {
				stats.DueReminders++
			}
		}

		results = append(results, stats)
	}

	c.JSON(200, results)
}

func (app *Application) getFuelStats(c *gin.Context) {
	vehicleID := c.Param("id")

	var fuelEntries []FuelEntry
	if err := app.db.Where("vehicle_id = ?", vehicleID).Order("date DESC").Find(&fuelEntries).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	type FuelStats struct {
		TotalCost     float64       `json:"total_cost"`
		AverageMPG    float64       `json:"average_mpg"`
		TotalGallons  float64       `json:"total_gallons"`
		TotalDistance float64       `json:"total_distance"`
		LastFillup    *FuelEntry    `json:"last_fillup"`
		MonthlyTrend  []gin.H       `json:"monthly_trend"`
	}

	stats := FuelStats{}

	if len(fuelEntries) == 0 {
		c.JSON(200, stats)
		return
	}

	stats.LastFillup = &fuelEntries[0]

	// Calculate totals
	firstOdo := fuelEntries[len(fuelEntries)-1].Odometer
	lastOdo := fuelEntries[0].Odometer

	for _, f := range fuelEntries {
		stats.TotalCost += f.Price
		stats.TotalGallons += f.Gallons
	}

	stats.TotalDistance = lastOdo - firstOdo
	if stats.TotalGallons > 0 {
		stats.AverageMPG = stats.TotalDistance / stats.TotalGallons
	}

	// Calculate monthly trend
	monthlyData := make(map[string]gin.H)
	for _, f := range fuelEntries {
		month := f.Date.Format("2006-01")
		if data, exists := monthlyData[month]; exists {
			data["cost"] = data["cost"].(float64) + f.Price
			data["gallons"] = data["gallons"].(float64) + f.Gallons
		} else {
			monthlyData[month] = gin.H{
				"month":   month,
				"cost":    f.Price,
				"gallons": f.Gallons,
			}
		}
	}

	for _, data := range monthlyData {
		stats.MonthlyTrend = append(stats.MonthlyTrend, data)
	}

	c.JSON(200, stats)
}

func (app *Application) getExpenseStats(c *gin.Context) {
	vehicleID := c.Param("id")

	var expenses []Expense
	if err := app.db.Where("vehicle_id = ?", vehicleID).Order("date DESC").Find(&expenses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	type CategoryStats struct {
		Category string  `json:"category"`
		Total    float64 `json:"total"`
		Count    int     `json:"count"`
	}

	categoryMap := make(map[string]CategoryStats)
	totalCost := 0.0

	for _, e := range expenses {
		totalCost += e.Amount
		if stats, exists := categoryMap[e.Category]; exists {
			stats.Total += e.Amount
			stats.Count++
			categoryMap[e.Category] = stats
		} else {
			categoryMap[e.Category] = CategoryStats{
				Category: e.Category,
				Total:    e.Amount,
				Count:    1,
			}
		}
	}

	var categories []CategoryStats
	for _, stats := range categoryMap {
		categories = append(categories, stats)
	}

	c.JSON(200, gin.H{
		"total_cost":    totalCost,
		"categories":    categories,
		"expense_count": len(expenses),
	})
}

// Fuel Entry Handlers with Enhanced Validation

func (app *Application) createFuelEntryEnhanced(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var req struct {
		Date      time.Time `json:"date" binding:"required"`
		Gallons   float64   `json:"gallons" binding:"required,gt=0"`
		Price     float64   `json:"price" binding:"required,gt=0"`
		Odometer  float64   `json:"odometer" binding:"required,gt=0"`
		Location  string    `json:"location"`
		Notes     string    `json:"notes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate odometer is greater than previous entry
	var lastEntry FuelEntry
	if err := app.db.Where("vehicle_id = ?", vehicleID).Order("odometer DESC").First(&lastEntry).Error; err == nil {
		if req.Odometer < lastEntry.Odometer {
			c.JSON(400, gin.H{"error": "Odometer cannot be less than previous entry"})
			return
		}
	}

	entry := FuelEntry{
		VehicleID: uint(vehicleID),
		Date:      req.Date,
		Gallons:   req.Gallons,
		Price:     req.Price,
		Odometer:  req.Odometer,
		Location:  req.Location,
		Notes:     req.Notes,
	}

	if err := app.db.Create(&entry).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := app.db.Model(&Vehicle{}).Where("id = ?", vehicleID).Update("odometer", req.Odometer).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update vehicle odometer"})
		return
	}

	// Check and trigger reminders
	alerts := app.checkVehicleReminders(uint(vehicleID), req.Odometer)

	c.JSON(201, gin.H{
		"entry":  entry,
		"alerts": alerts,
	})
}

// Expense Handlers

func (app *Application) createExpenseEnhanced(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var req struct {
		Category string    `json:"category" binding:"required"`
		Amount   float64   `json:"amount" binding:"required,gt=0"`
		Date     time.Time `json:"date" binding:"required"`
		Notes    string    `json:"notes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	expense := Expense{
		VehicleID: uint(vehicleID),
		Category:  req.Category,
		Amount:    req.Amount,
		Date:      req.Date,
		Notes:     req.Notes,
	}

	if err := app.db.Create(&expense).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, expense)
}

// Reminder Handlers with Advanced Logic

func (app *Application) checkRemindersDue(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var vehicle Vehicle
	if err := app.db.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Vehicle not found"})
		return
	}

	alerts := app.checkVehicleReminders(uint(vehicleID), vehicle.Odometer)
	c.JSON(200, gin.H{
		"vehicle": vehicle,
		"alerts":  alerts,
	})
}

func (app *Application) completeReminder(c *gin.Context) {
	reminderID := c.Param("id")

	var req struct {
		ServiceDate  time.Time `json:"service_date" binding:"required"`
		ServiceMiles float64   `json:"service_miles" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := app.db.Model(&MaintenanceReminder{}).Where("id = ?", reminderID).Updates(map[string]interface{}{
		"last_service_date":  req.ServiceDate,
		"last_service_miles": req.ServiceMiles,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Reminder completed"})
}

func (app *Application) getRemindersOverdue(c *gin.Context) {
	userID := c.GetUint("userID")

	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	type ReminderAlert struct {
		VehicleID    uint    `json:"vehicle_id"`
		VehicleName  string  `json:"vehicle_name"`
		ReminderID   uint    `json:"reminder_id"`
		ReminderName string  `json:"reminder_name"`
		Status       string  `json:"status"` // overdue, soon, upcoming
		MilesToGo    float64 `json:"miles_to_go"`
		DaysUntil    int     `json:"days_until"`
	}

	var allAlerts []ReminderAlert

	for _, v := range vehicles {
		var reminders []MaintenanceReminder
		app.db.Where("vehicle_id = ?", v.ID).Find(&reminders)

		for _, r := range reminders {
			alert := ReminderAlert{
				VehicleID:    v.ID,
				VehicleName:  fmt.Sprintf("%d %s %s", v.Year, v.Make, v.Model),
				ReminderID:   r.ID,
				ReminderName: r.Name,
				Status:       "upcoming",
			}

			// Check mileage
			if r.IntervalMiles > 0 {
				nextServiceMiles := r.LastServiceMiles + r.IntervalMiles
				alert.MilesToGo = nextServiceMiles - v.Odometer

				if v.Odometer >= nextServiceMiles {
					alert.Status = "overdue"
				} else if alert.MilesToGo < 500 {
					alert.Status = "soon"
				}
			}

			// Check days
			if r.IntervalDays > 0 {
				nextServiceDate := r.LastServiceDate.AddDate(0, 0, r.IntervalDays)
				alert.DaysUntil = int(time.Until(nextServiceDate).Hours() / 24)

				if alert.DaysUntil <= 0 {
					alert.Status = "overdue"
				} else if alert.DaysUntil < 7 {
					alert.Status = "soon"
				}
			}

			if alert.Status != "upcoming" {
				allAlerts = append(allAlerts, alert)
			}
		}
	}

	c.JSON(200, gin.H{"alerts": allAlerts})
}

// Attachment Handlers

func (app *Application) attachFileToEntry(c *gin.Context) {
	entryType := c.Query("type") // fuel or expense
	entryID := c.Query("entry_id")

	if entryType == "" || entryID == "" {
		c.JSON(400, gin.H{"error": "Missing parameters"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(400, gin.H{"error": "File too large (max 5MB)"})
		return
	}

	assetsPath := "/assets"
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filepath := assetsPath + "/" + filename

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Upload failed"})
		return
	}

	// Create attachment record
	attachment := Attachment{
		EntryID:   parseUint(entryID),
		EntryType: entryType,
		Filename:  file.Filename,
		Path:      filepath,
	}

	if err := app.db.Create(&attachment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(201, attachment)
}

// Sharing Handlers

func (app *Application) listVehicleUsers(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicleUsers []VehicleUser
	if err := app.db.Where("vehicle_id = ?", vehicleID).Find(&vehicleUsers).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var users []User
	for _, vu := range vehicleUsers {
		var user User
		app.db.First(&user, vu.UserID)
		users = append(users, user)
	}

	c.JSON(200, users)
}

func (app *Application) removeVehicleUser(c *gin.Context) {
	vehicleID := c.Param("id")
	userID := c.Param("userId")

	if err := app.db.Where("vehicle_id = ? AND user_id = ?", vehicleID, userID).Delete(&VehicleUser{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User removed from vehicle"})
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

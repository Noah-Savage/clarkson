package main

import (
	"os"
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

type Application struct {
	db *gorm.DB
	router *gin.Engine
	jwtSecret string
}

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Get JWT secret from env or generate
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-secret-change-in-production"
	}

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create app instance
	app := &Application{
		db:        db,
		router:    router,
		jwtSecret: jwtSecret,
	}

	// Setup routes
	setupRoutes(app)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Clarkson starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}

func initDB() (*gorm.DB, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/config"
	}

	// Ensure config directory exists
	os.MkdirAll(configPath, 0755)

	dbPath := configPath + "/clarkson.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate models
	return db, db.AutoMigrate(
		&User{},
		&Vehicle{},
		&VehicleUser{},
		&FuelEntry{},
		&Expense{},
		&MaintenanceReminder{},
		&Attachment{},
		&Notification{},
	)
}

func setupRoutes(app *Application) {
	// Auth routes
	auth := app.router.Group("/api/auth")
	{
		auth.POST("/register", app.handleRegister)
		auth.POST("/login", app.handleLogin)
	}

	// Protected routes
	protected := app.router.Group("/api")
	protected.Use(authMiddleware(app.jwtSecret))
	{
		// User routes
		protected.GET("/users/:id", app.getUser)
		protected.PUT("/users/:id", app.updateUser)

		// Vehicle routes
		protected.GET("/vehicles", app.listVehicles)
		protected.POST("/vehicles", app.createVehicle)
		protected.GET("/vehicles/:id", app.getVehicle)
		protected.PUT("/vehicles/:id", app.updateVehicle)
		protected.DELETE("/vehicles/:id", app.deleteVehicle)
		protected.POST("/vehicles/:id/share", app.shareVehicle)

		// Fuel entry routes
		protected.GET("/vehicles/:id/fuel", app.listFuelEntries)
		protected.POST("/vehicles/:id/fuel", app.createFuelEntry)
		protected.PUT("/fuel/:id", app.updateFuelEntry)
		protected.DELETE("/fuel/:id", app.deleteFuelEntry)

		// Expense routes
		protected.GET("/vehicles/:id/expenses", app.listExpenses)
		protected.POST("/vehicles/:id/expenses", app.createExpense)
		protected.PUT("/expenses/:id", app.updateExpense)
		protected.DELETE("/expenses/:id", app.deleteExpense)

		// Reminder routes
		protected.GET("/vehicles/:id/reminders", app.listReminders)
		protected.POST("/vehicles/:id/reminders", app.createReminder)
		protected.PUT("/reminders/:id", app.updateReminder)
		protected.DELETE("/reminders/:id", app.deleteReminder)
		protected.GET("/reminders/check", app.checkReminders)

		// Reports
		protected.GET("/vehicles/:id/report", app.generateReport)
		protected.GET("/report/overall", app.generateOverallReport)
		protected.GET("/export/csv", app.exportCSV)
		protected.GET("/export/pdf", app.exportPDF)

		// Import/Migration
		protected.POST("/import/hammond", app.importHammond)
		protected.POST("/import/fuelly", app.importFuelly)
		protected.POST("/import/clarkson", app.importClarkson)

		// Notification routes
		protected.GET("/notifications", app.listNotifications)
		protected.POST("/notifications", app.createNotification)
		protected.PUT("/notifications/:id", app.updateNotification)
		protected.DELETE("/notifications/:id", app.deleteNotification)
	}

	// File upload/download (protected)
	protected.POST("/upload", app.uploadFile)
	protected.GET("/download/:id", app.downloadFile)
}

// JWT Claims
type Claims struct {
	UserID uint
	Email  string
	jwt.RegisteredClaims
}

func (app *Application) handleRegister(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user := &User{
		Email:    req.Email,
		Name:     req.Name,
		Role:     "admin",
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(500, gin.H{"error": "Failed to process password"})
		return
	}

	result := app.db.Create(user)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(201, user)
}

func (app *Application) handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := app.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(500, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func (app *Application) getUser(c *gin.Context) {
	userID := c.GetUint("userID")
	var user User
	if err := app.db.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func (app *Application) updateUser(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
		Units    string `json:"units"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := app.db.Model(&User{}).Where("id = ?", userID).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) listVehicles(c *gin.Context) {
	userID := c.GetUint("userID")
	var vehicles []Vehicle
	if err := app.db.
		Joins("LEFT JOIN vehicle_users ON vehicle_users.vehicle_id = vehicles.id").
		Where("vehicles.user_id = ? OR vehicle_users.user_id = ?", userID, userID).
		Find(&vehicles).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, vehicles)
}

func (app *Application) createVehicle(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Make         string  `json:"make" binding:"required"`
		Model        string  `json:"model" binding:"required"`
		Year         int     `json:"year" binding:"required"`
		Odometer     float64 `json:"odometer"`
		MileageUnit  string  `json:"mileage_unit"` // km or mi
		FuelType     string  `json:"fuel_type"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	vehicle := Vehicle{
		UserID:      userID,
		Make:        req.Make,
		Model:       req.Model,
		Year:        req.Year,
		Odometer:    req.Odometer,
		MileageUnit: req.MileageUnit,
		FuelType:    req.FuelType,
	}

	if err := app.db.Create(&vehicle).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, vehicle)
}

func (app *Application) getVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	var vehicle Vehicle
	if err := app.db.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Vehicle not found"})
		return
	}
	c.JSON(200, vehicle)
}

func (app *Application) updateVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	var req struct {
		Make        string  `json:"make"`
		Model       string  `json:"model"`
		Year        int     `json:"year"`
		Odometer    float64 `json:"odometer"`
		MileageUnit string  `json:"mileage_unit"`
		FuelType    string  `json:"fuel_type"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := app.db.Model(&Vehicle{}).Where("id = ?", vehicleID).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) deleteVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	if err := app.db.Delete(&Vehicle{}, vehicleID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}

func (app *Application) shareVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := app.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	vu := VehicleUser{
		VehicleID: parseUint(vehicleID),
		UserID:    user.ID,
	}

	if err := app.db.Create(&vu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, vu)
}

func (app *Application) listFuelEntries(c *gin.Context) {
	vehicleID := c.Param("id")
	var entries []FuelEntry
	if err := app.db.Where("vehicle_id = ?", vehicleID).Order("date DESC").Find(&entries).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, entries)
}

func (app *Application) createFuelEntry(c *gin.Context) {
	vehicleID := c.Param("id")
	var req struct {
		Date      time.Time `json:"date" binding:"required"`
		Gallons   float64   `json:"gallons" binding:"required"`
		Price     float64   `json:"price" binding:"required"`
		Odometer  float64   `json:"odometer" binding:"required"`
		Location  string    `json:"location"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	entry := FuelEntry{
		VehicleID: parseUint(vehicleID),
		Date:      req.Date,
		Gallons:   req.Gallons,
		Price:     req.Price,
		Odometer:  req.Odometer,
		Location:  req.Location,
	}

	if err := app.db.Create(&entry).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check reminders for this vehicle
	app.checkVehicleReminders(parseUint(vehicleID), req.Odometer)

	c.JSON(201, entry)
}

func (app *Application) updateFuelEntry(c *gin.Context) {
	fuelID := c.Param("id")
	var req struct {
		Date      time.Time `json:"date"`
		Gallons   float64   `json:"gallons"`
		Price     float64   `json:"price"`
		Odometer  float64   `json:"odometer"`
		Location  string    `json:"location"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var entry FuelEntry
	if err := app.db.First(&entry, fuelID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Entry not found"})
		return
	}

	if err := app.db.Model(&FuelEntry{}).Where("id = ?", fuelID).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) deleteFuelEntry(c *gin.Context) {
	fuelID := c.Param("id")
	if err := app.db.Delete(&FuelEntry{}, fuelID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}

func (app *Application) listExpenses(c *gin.Context) {
	vehicleID := c.Param("id")
	var expenses []Expense
	if err := app.db.Where("vehicle_id = ?", vehicleID).Order("date DESC").Find(&expenses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, expenses)
}

func (app *Application) createExpense(c *gin.Context) {
	vehicleID := c.Param("id")
	var req struct {
		Category string    `json:"category" binding:"required"`
		Amount   float64   `json:"amount" binding:"required"`
		Date     time.Time `json:"date" binding:"required"`
		Notes    string    `json:"notes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	expense := Expense{
		VehicleID: parseUint(vehicleID),
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

func (app *Application) updateExpense(c *gin.Context) {
	expenseID := c.Param("id")
	var req struct {
		Category string    `json:"category"`
		Amount   float64   `json:"amount"`
		Date     time.Time `json:"date"`
		Notes    string    `json:"notes"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := app.db.Model(&Expense{}).Where("id = ?", expenseID).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) deleteExpense(c *gin.Context) {
	expenseID := c.Param("id")
	if err := app.db.Delete(&Expense{}, expenseID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}

func (app *Application) listReminders(c *gin.Context) {
	vehicleID := c.Param("id")
	var reminders []MaintenanceReminder
	if err := app.db.Where("vehicle_id = ?", vehicleID).Find(&reminders).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, reminders)
}

func (app *Application) createReminder(c *gin.Context) {
	vehicleID := c.Param("id")
	var req struct {
		Name             string    `json:"name" binding:"required"`
		IntervalMiles    float64   `json:"interval_miles"`
		IntervalDays     int       `json:"interval_days"`
		LastServiceDate  time.Time `json:"last_service_date"`
		LastServiceMiles float64   `json:"last_service_miles"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	reminder := MaintenanceReminder{
		VehicleID:        parseUint(vehicleID),
		Name:             req.Name,
		IntervalMiles:    req.IntervalMiles,
		IntervalDays:     req.IntervalDays,
		LastServiceDate:  req.LastServiceDate,
		LastServiceMiles: req.LastServiceMiles,
	}

	if err := app.db.Create(&reminder).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, reminder)
}

func (app *Application) updateReminder(c *gin.Context) {
	reminderID := c.Param("id")
	var req struct {
		Name             string    `json:"name"`
		IntervalMiles    float64   `json:"interval_miles"`
		IntervalDays     int       `json:"interval_days"`
		LastServiceDate  time.Time `json:"last_service_date"`
		LastServiceMiles float64   `json:"last_service_miles"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := app.db.Model(&MaintenanceReminder{}).Where("id = ?", reminderID).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) deleteReminder(c *gin.Context) {
	reminderID := c.Param("id")
	if err := app.db.Delete(&MaintenanceReminder{}, reminderID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}

func (app *Application) checkReminders(c *gin.Context) {
	userID := c.GetUint("userID")
	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	var alerts []gin.H
	for _, v := range vehicles {
		reminders := app.checkVehicleReminders(v.ID, v.Odometer)
		alerts = append(alerts, reminders...)
	}

	c.JSON(200, gin.H{"alerts": alerts})
}

func (app *Application) checkVehicleReminders(vehicleID uint, currentOdometer float64) []gin.H {
	var reminders []MaintenanceReminder
	app.db.Where("vehicle_id = ?", vehicleID).Find(&reminders)

	var alerts []gin.H

	for _, r := range reminders {
		var isOverdue bool
		var daysUntilDue int
		var milesToGo float64

		// Check mileage-based reminder
		if r.IntervalMiles > 0 {
			nextServiceMiles := r.LastServiceMiles + r.IntervalMiles
			if currentOdometer >= nextServiceMiles {
				isOverdue = true
				milesToGo = nextServiceMiles - currentOdometer
			} else {
				milesToGo = nextServiceMiles - currentOdometer
			}
		}

		// Check days-based reminder
		if r.IntervalDays > 0 {
			nextServiceDate := r.LastServiceDate.AddDate(0, 0, r.IntervalDays)
			daysUntilDue = int(time.Until(nextServiceDate).Hours() / 24)
			if daysUntilDue <= 0 {
				isOverdue = true
			}
		}

		status := "upcoming"
		if isOverdue {
			status = "overdue"
		} else if daysUntilDue < 7 || milesToGo < 500 {
			status = "soon"
		}

		alert := gin.H{
			"vehicleID":    vehicleID,
			"reminderID":   r.ID,
			"name":         r.Name,
			"status":       status,
			"daysUntilDue": daysUntilDue,
			"milesToGo":    milesToGo,
		}

		if status != "upcoming" {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

func (app *Application) generateReport(c *gin.Context) {
	vehicleID := c.Param("id")

	var vehicle Vehicle
	if err := app.db.First(&vehicle, vehicleID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Vehicle not found"})
		return
	}

	var fuelEntries []FuelEntry
	app.db.Where("vehicle_id = ?", vehicleID).Order("date").Find(&fuelEntries)

	var expenses []Expense
	app.db.Where("vehicle_id = ?", vehicleID).Find(&expenses)

	totalCost := 0.0
	totalGallons := 0.0
	avgMPG := 0.0

	for _, f := range fuelEntries {
		totalCost += f.Price
		totalGallons += f.Gallons
	}

	for _, e := range expenses {
		totalCost += e.Amount
	}

	if totalGallons > 0 {
		avgMPG = (vehicle.Odometer - (fuelEntries[0].Odometer)) / totalGallons
	}

	c.JSON(200, gin.H{
		"vehicle":      vehicle,
		"fuelEntries":  fuelEntries,
		"expenses":     expenses,
		"totalCost":    totalCost,
		"totalGallons": totalGallons,
		"avgMPG":       avgMPG,
	})
}

func (app *Application) generateOverallReport(c *gin.Context) {
	userID := c.GetUint("userID")
	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	type VehicleStats struct {
		Vehicle    Vehicle
		TotalCost  float64
		TotalMiles float64
		AvgMPG     float64
	}

	var stats []VehicleStats
	totalCost := 0.0

	for _, v := range vehicles {
		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Order("date").Find(&fuelEntries)

		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Find(&expenses)

		vCost := 0.0
		vGallons := 0.0

		for _, f := range fuelEntries {
			vCost += f.Price
		}

		for _, e := range expenses {
			vCost += e.Amount
		}

		for _, f := range fuelEntries {
			vGallons += f.Gallons
		}

		totalCost += vCost

		stats = append(stats, VehicleStats{
			Vehicle:   v,
			TotalCost: vCost,
		})
	}

	c.JSON(200, gin.H{
		"vehicles":  stats,
		"totalCost": totalCost,
	})
}

func (app *Application) exportCSV(c *gin.Context) {
	userID := c.GetUint("userID")
	c.Header("Content-Disposition", "attachment; filename=clarkson-export.csv")
	c.Header("Content-Type", "text/csv")

	csv := "Vehicle,Date,Type,Amount,Odometer,Notes\n"

	var vehicles []Vehicle
	app.db.Where("user_id = ?", userID).Find(&vehicles)

	for _, v := range vehicles {
		var fuelEntries []FuelEntry
		app.db.Where("vehicle_id = ?", v.ID).Find(&fuelEntries)

		for _, f := range fuelEntries {
			csv += fmt.Sprintf("%s %s,%s,Fuel,%.2f,%.1f,\n", v.Year, v.Make, f.Date.Format("2006-01-02"), f.Price, f.Odometer)
		}

		var expenses []Expense
		app.db.Where("vehicle_id = ?", v.ID).Find(&expenses)

		for _, e := range expenses {
			csv += fmt.Sprintf("%s %s,%s,%s,%.2f,,\"%s\"\n", v.Year, v.Make, e.Date.Format("2006-01-02"), e.Category, e.Amount, e.Notes)
		}
	}

	c.Data(200, "text/csv", []byte(csv))
}

func (app *Application) exportPDF(c *gin.Context) {
	c.JSON(200, gin.H{"message": "PDF export not yet implemented"})
}

func (app *Application) importHammond(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hammond import not yet implemented"})
}

func (app *Application) importFuelly(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Fuelly import not yet implemented"})
}

func (app *Application) importClarkson(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Clarkson import not yet implemented"})
}

func (app *Application) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath = "/assets"
	}

	os.MkdirAll(assetsPath, 0755)

	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	filepath := assetsPath + "/" + filename

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(201, gin.H{"filename": filename})
}

func (app *Application) downloadFile(c *gin.Context) {
	filename := c.Param("id")
	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath = "/assets"
	}

	filepath := assetsPath + "/" + filename
	if _, err := os.Stat(filepath); err != nil {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.File(filepath)
}

func (app *Application) listNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	var notifications []Notification
	if err := app.db.Where("user_id = ?", userID).Find(&notifications).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, notifications)
}

func (app *Application) createNotification(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	notification := Notification{
		UserID:  userID,
		Message: req.Message,
	}

	if err := app.db.Create(&notification).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, notification)
}

func (app *Application) updateNotification(c *gin.Context) {
	notificationID := c.Param("id")
	var req struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var notification Notification
	if err := app.db.First(&notification, notificationID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Notification not found"})
		return
	}

	if err := app.db.Model(&notification).Updates(req).Error; err != nil {
		c.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Updated"})
}

func (app *Application) deleteNotification(c *gin.Context) {
	notificationID := c.Param("id")
	if err := app.db.Delete(&Notification{}, notificationID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(200, gin.H{"message": "Deleted"})
}

func authMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "No token"})
			c.Abort()
			return
		}

		claims := &Claims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func parseUint(s string) uint {
	var u uint
	fmt.Sscanf(s, "%d", &u)
	return u
}

type User struct {
	ID       uint
	Email    string
	Name     string
	Password string
	Role     string
}

type Vehicle struct {
	ID          uint
	UserID      uint
	Make        string
	Model       string
	Year        int
	Odometer    float64
	MileageUnit string
	FuelType    string
}

type VehicleUser struct {
	VehicleID uint
	UserID    uint
}

type FuelEntry struct {
	ID        uint
	VehicleID uint
	Date      time.Time
	Gallons   float64
	Price     float64
	Odometer  float64
	Location  string
}

type Expense struct {
	ID        uint
	VehicleID uint
	Category  string
	Amount    float64
	Date      time.Time
	Notes     string
}

type MaintenanceReminder struct {
	ID               uint
	VehicleID        uint
	Name             string
	IntervalMiles    float64
	IntervalDays     int
	LastServiceDate  time.Time
	LastServiceMiles float64
}

type Attachment struct {
	ID       uint
	VehicleID uint
	Filename string
}

type Notification struct {
	ID       uint
	UserID   uint
	Message  string
}

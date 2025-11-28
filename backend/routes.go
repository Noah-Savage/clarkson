package main

import (
	"github.com/gin-gonic/gin"
)


func setupRoutes(app *Application) {
	// Auth routes (no auth required)
	auth := app.router.Group("/api/auth")
	{
		auth.POST("/register", app.handleRegister)
		auth.POST("/login", app.handleLogin)
	}

	// Protected routes (require auth)
	protected := app.router.Group("/api")
	protected.Use(authMiddleware(app.jwtSecret))
	{
		// User routes
		protected.GET("/users/:id", app.getUser)
		protected.PUT("/users/:id", app.updateUser)

		// Vehicle routes - enhanced
		protected.GET("/vehicles", app.listVehiclesWithStats)
		protected.POST("/vehicles", app.createVehicle)
		protected.GET("/vehicles/:id", app.getVehicle)
		protected.PUT("/vehicles/:id", app.updateVehicle)
		protected.DELETE("/vehicles/:id", app.deleteVehicle)
		protected.POST("/vehicles/:id/share", app.shareVehicle)
		protected.GET("/vehicles/:id/users", app.listVehicleUsers)
		protected.DELETE("/vehicles/:id/users/:userId", app.removeVehicleUser)

		// Fuel entry routes - enhanced
		protected.GET("/vehicles/:id/fuel", app.listFuelEntries)
		protected.POST("/vehicles/:id/fuel", app.createFuelEntryEnhanced)
		protected.GET("/vehicles/:id/fuel-stats", app.getFuelStats)
		protected.PUT("/fuel/:id", app.updateFuelEntry)
		protected.DELETE("/fuel/:id", app.deleteFuelEntry)

		// Expense routes - enhanced
		protected.GET("/vehicles/:id/expenses", app.listExpenses)
		protected.POST("/vehicles/:id/expenses", app.createExpenseEnhanced)
		protected.GET("/vehicles/:id/expense-stats", app.getExpenseStats)
		protected.PUT("/expenses/:id", app.updateExpense)
		protected.DELETE("/expenses/:id", app.deleteExpense)

		// Reminder routes - enhanced
		protected.GET("/vehicles/:id/reminders", app.listReminders)
		protected.POST("/vehicles/:id/reminders", app.createReminder)
		protected.PUT("/reminders/:id", app.updateReminder)
		protected.DELETE("/reminders/:id", app.deleteReminder)
		protected.POST("/reminders/:id/complete", app.completeReminder)
		protected.GET("/reminders/check", app.checkReminders)
		protected.GET("/reminders/overdue", app.getRemindersOverdue)
		protected.GET("/vehicles/:id/reminders/due", app.checkRemindersDue)

		// Notification routes
		protected.GET("/notifications", app.getUnreadNotifications)
		protected.GET("/notifications/summary", app.getNotificationSummary)
		protected.POST("/notifications/:id/read", app.markNotificationRead)
		protected.POST("/notifications/:id/dismiss", app.dismissNotification)

		// Reports
		protected.GET("/vehicles/:id/report", app.generateReport)
		protected.GET("/report/overall", app.generateOverallReport)
		protected.GET("/export/csv", app.exportCSV)
		protected.GET("/export/pdf", app.exportPDF)

		// Import/Migration
		protected.POST("/import/hammond", app.importHammond)
		protected.POST("/import/fuelly", app.importFuelly)
		protected.POST("/import/clarkson", app.importClarkson)

		// File upload/download
		protected.POST("/upload", app.uploadFile)
		protected.GET("/download/:id", app.downloadFile)
		protected.POST("/attach", app.attachFileToEntry)
	}

	// Health check (no auth)
	app.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}

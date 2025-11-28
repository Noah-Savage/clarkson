package main

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)


type Notification struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	VehicleID     uint      `json:"vehicle_id"`
	ReminderID    uint      `json:"reminder_id"`
	Type          string    `json:"type"` // reminder_due, reminder_overdue
	Title         string    `json:"title"`
	Message       string    `json:"message"`
	Status        string    `json:"status"` // unread, read, dismissed
	CreatedAt     time.Time `json:"created_at"`
	DismissedAt   *time.Time `json:"dismissed_at"`
}

func (app *Application) checkVehicleRemindersAdvanced(vehicleID uint, currentOdometer float64) []Notification {
	var reminders []MaintenanceReminder
	app.db.Where("vehicle_id = ?", vehicleID).Find(&reminders)

	var notifications []Notification

	for _, r := range reminders {
		var notif Notification
		notif.VehicleID = vehicleID
		notif.ReminderID = r.ID
		notif.Status = "unread"
		notif.CreatedAt = time.Now()

		// Check mileage-based reminder
		if r.IntervalMiles > 0 {
			nextServiceMiles := r.LastServiceMiles + r.IntervalMiles
			milesSinceService := nextServiceMiles - currentOdometer

			if currentOdometer >= nextServiceMiles {
				notif.Type = "reminder_overdue"
				notif.Title = fmt.Sprintf("%s - OVERDUE", r.Name)
				notif.Message = fmt.Sprintf("This service was due %.0f miles ago at %.0f miles", -milesSinceService, nextServiceMiles)
				notifications = append(notifications, notif)
			} else if milesSinceService < 500 {
				notif.Type = "reminder_due"
				notif.Title = fmt.Sprintf("%s - DUE SOON", r.Name)
				notif.Message = fmt.Sprintf("Service due in %.0f miles (at %.0f miles)", milesSinceService, nextServiceMiles)
				notifications = append(notifications, notif)
			}
		}

		// Check days-based reminder
		if r.IntervalDays > 0 {
			nextServiceDate := r.LastServiceDate.AddDate(0, 0, r.IntervalDays)
			daysUntilDue := int(time.Until(nextServiceDate).Hours() / 24)

			if daysUntilDue <= 0 {
				notif.Type = "reminder_overdue"
				notif.Title = fmt.Sprintf("%s - OVERDUE", r.Name)
				notif.Message = fmt.Sprintf("This service was due %d days ago", -daysUntilDue)
				notifications = append(notifications, notif)
			} else if daysUntilDue < 7 {
				notif.Type = "reminder_due"
				notif.Title = fmt.Sprintf("%s - DUE SOON", r.Name)
				notif.Message = fmt.Sprintf("Service due in %d days", daysUntilDue)
				notifications = append(notifications, notif)
			}
		}
	}

	return notifications
}

func (app *Application) storeNotifications(userID uint, notifications []Notification) error {
	for i := range notifications {
		notifications[i].UserID = userID
	}
	return app.db.Create(&notifications).Error
}

func (app *Application) getUnreadNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	var notifications []Notification

	if err := app.db.
		Where("user_id = ? AND status = ?", userID, "unread").
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, notifications)
}

func (app *Application) markNotificationRead(c *gin.Context) {
	notifID := c.Param("id")

	if err := app.db.Model(&Notification{}).Where("id = ?", notifID).Update("status", "read").Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Marked as read"})
}

func (app *Application) dismissNotification(c *gin.Context) {
	notifID := c.Param("id")
	now := time.Now()

	if err := app.db.Model(&Notification{}).Where("id = ?", notifID).Updates(map[string]interface{}{
		"status":       "dismissed",
		"dismissed_at": now,
	}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Dismissed"})
}

func (app *Application) getNotificationSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	var unreadCount int64
	app.db.Model(&Notification{}).Where("user_id = ? AND status = ?", userID, "unread").Count(&unreadCount)

	var overdueCount int64
	app.db.Model(&Notification{}).Where("user_id = ? AND status != ? AND type = ?", userID, "dismissed", "reminder_overdue").Count(&overdueCount)

	var upcomingCount int64
	app.db.Model(&Notification{}).Where("user_id = ? AND status != ? AND type = ?", userID, "dismissed", "reminder_due").Count(&upcomingCount)

	c.JSON(200, gin.H{
		"unread_count":   unreadCount,
		"overdue_count":  overdueCount,
		"upcoming_count": upcomingCount,
	})
}

package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
	"gorm.io/gorm"
)


type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Name      string    json:"name"`
	Password  string    json:"-"`
	Role      string    json:"role"` // admin, user
	Currency  string    `json:"currency"` // USD, EUR, etc
	Units     string    `json:"units"` // mi, km
	CreatedAt time.Time json:"created_at"`
	UpdatedAt time.Time json:"updated_at"`

	Vehicles []Vehicle `gorm:"foreignKey:UserID" json:"vehicles,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

type Vehicle struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"` // Owner
	Make        string    `json:"make"`
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	Odometer    float64   `json:"odometer"` // Current odometer reading
	MileageUnit string    `json:"mileage_unit"` // mi or km
	FuelType    string    `json:"fuel_type"` // Petrol, Diesel, Electric, Hybrid, etc
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	FuelEntries []FuelEntry           `gorm:"foreignKey:VehicleID" json:"fuel_entries,omitempty"`
	Expenses    []Expense             `gorm:"foreignKey:VehicleID" json:"expenses,omitempty"`
	Reminders   []MaintenanceReminder `gorm:"foreignKey:VehicleID" json:"reminders,omitempty"`
	SharedUsers []VehicleUser         `gorm:"foreignKey:VehicleID" json:"shared_users,omitempty"`
}

type VehicleUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VehicleID uint      `json:"vehicle_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FuelEntry struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VehicleID uint      `json:"vehicle_id"`
	Date      time.Time `json:"date"`
	Gallons   float64   `json:"gallons"` // Or liters
	Price     float64   `json:"price"`
	Odometer  float64   `json:"odometer"`
	Location  string    `json:"location"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Vehicle     Vehicle      `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:EntryID;foreignKeyValue:fuelentry" json:"attachments,omitempty"`
}

type Expense struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VehicleID uint      `json:"vehicle_id"`
	Category  string    `json:"category"` // Maintenance, Insurance, Parking, etc
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Vehicle     Vehicle      `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:EntryID;foreignKeyValue:expense" json:"attachments,omitempty"`
}

type MaintenanceReminder struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	VehicleID        uint      `json:"vehicle_id"`
	Name             string    `json:"name"` // Oil Change, Tire Rotation, etc
	IntervalMiles    float64   `json:"interval_miles"` // 0 = disabled
	IntervalDays     int       `json:"interval_days"` // 0 = disabled
	LastServiceDate  time.Time `json:"last_service_date"`
	LastServiceMiles float64   `json:"last_service_miles"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Vehicle Vehicle `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
}

type Attachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EntryID   uint      `json:"entry_id"`
	EntryType string    `json:"entry_type"` // fuelentry, expense
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `json:"user_id"`
	VehicleID     uint       `json:"vehicle_id"`
	ReminderID    uint       `json:"reminder_id"`
	Type          string     `json:"type"` // reminder_due, reminder_overdue
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	Status        string     `json:"status"` // unread, read, dismissed
	CreatedAt     time.Time  `json:"created_at"`
	DismissedAt   *time.Time `json:"dismissed_at"`
}

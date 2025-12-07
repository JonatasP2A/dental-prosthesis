package client

import "time"

// Client represents a dental clinic or dentist who orders prosthetic work
type Client struct {
	ID          string
	LaboratoryID string
	Name        string
	Email       string
	Phone       string
	Address     Address
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Address represents a client's address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}


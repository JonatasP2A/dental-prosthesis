package order

import "time"

// Order represents a prosthesis work order from a client
type Order struct {
	ID          string
	ClientID    string
	LaboratoryID string
	Status      Status
	Prosthesis  []ProsthesisItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Status represents the workflow state of an order
type Status string

const (
	StatusReceived    Status = "received"
	StatusInProduction Status = "in_production"
	StatusQualityCheck Status = "quality_check"
	StatusReady       Status = "ready"
	StatusDelivered   Status = "delivered"
	StatusRevision    Status = "revision"
)

// ProsthesisItem represents a prosthesis item in an order
type ProsthesisItem struct {
	Type       string
	Material   string
	Shade      string
	Quantity   int
	Notes      string
}


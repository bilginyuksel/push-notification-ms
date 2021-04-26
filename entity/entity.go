package entity

import "time"

type Application struct {
	// RecordID it is a unique id for application
	RecordID string

	Name        string
	Description string

	// Stores applications status - Active, Approving, Cancelled, Suspended
	Status string

	CreateTime time.Time

	// CancelTime when application is canceled it will be recorded.
	// Value is nil if it is not canceled.
	CancelTime *time.Time
}

type Client struct {
	// RecordID is a unique for any client registered to the app
	RecordID string

	// ClientID is an identifier for the client device
	ClientID string

	// APPID the app whom client registered
	APPID string

	Status string

	RegisterTime time.Time

	LastStatusChangeTime time.Time
}

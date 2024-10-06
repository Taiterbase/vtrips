package models

import (
	"time"

	validate "github.com/go-playground/validator/v10"
	ulid "github.com/oklog/ulid/v2"
)

type Trip interface {
	GetID() string
	GetClientID() string
	GetCreatedAt() int64
	GetUpdatedAt() int64
	GetTripType() TripType
	GetDeletedAt() int64

	SetID(string)
	SetClientID(string)
	SetCreatedAt(int64)
	SetUpdatedAt(int64)
	SetTripType(TripType)
	SetDeletedAt(int64)

	Validate() error
}

type TripStatus string

const (
	TripStatusDraft    TripStatus = "draft"
	TripStatusComplete TripStatus = "complete"
	TripStatusListed   TripStatus = "listed"
	TripStatusUnlisted TripStatus = "unlisted"
	TripStatusArchived TripStatus = "archived"
)

type TripBase struct {
	PK        string `json:"-" dynamodbav:"pk"`
	SK        string `json:"-" dynamodbav:"sk"`
	ID        string `json:"id" dynamodbav:"id" validate:"required"`
	ClientID  string `json:"client_id" dynamodbav:"client_id" validate:"required"`
	CreatedAt int64  `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt int64  `json:"updated_at" dynamodbav:"updated_at" updateable:"true"`
	DeletedAt int64  `json:"deleted_at" dynamodbav:"deleted_at" updateable:"true"`

	HousingType    HousingType `json:"housing_type" dynamodbav:"housing_type" updateable:"true"`
	PrivacyType    PrivacyType `json:"privacy_type" dynamodbav:"privacy_type" updateable:"true"`
	TripType       TripType    `json:"trip_type" dynamodbav:"trip_type" updateable:"true"`
	Latitude       float64     `json:"latitude" dynamodbav:"latitude" updateable:"true"`
	Longtitude     float64     `json:"longtitude" dynamodbav:"longtitude" updateable:"true"`
	Status         TripStatus  `json:"status" dynamodbav:"status" updateable:"true"`
	VolunteerLimit int         `json:"volunteer_limit" dynamodbav:"volunteer_limit" updateable:"true"`
	Name           string      `json:"name" dynamodbav:"name" updateable:"true"`
	Description    string      `json:"description" dynamodbav:"description" updateable:"true"`
	Mission        string      `json:"mission" dynamodbav:"mission" updateable:"true"`
	Price          float64     `json:"price" dynamodbav:"price" updateable:"true"`
	Currency       string      `json:"currency" dynamodbav:"currency" updateable:"true"`
	StartDate      int64       `json:"start_date" dynamodbav:"start_date" updateable:"true"`
	EndDate        int64       `json:"end_date" dynamodbav:"end_date" updateable:"true"`
}

var updateFields = []string{
	"updated_at",
	"deleted_at",
	"housing_type",
	"privacy_type",
	"trip_type",
	"latitude",
	"longtitude",
	"status",
	"volunteer_limit",
	"name",
	"description",
	"mission",
	"price",
	"currency",
	"start_date",
	"end_date",
}

func NewTrip() Trip {
	now := time.Now().UnixMilli()
	return &TripBase{
		ID:          ulid.Make().String(),
		CreatedAt:   now,
		UpdatedAt:   now,
		HousingType: OtherHousing,
		PrivacyType: OtherPrivacy,
		TripType:    OtherTrip,
		Status:      TripStatusDraft,
	}
}

func (t *TripBase) Validate() error {
	validate := validate.New()
	return validate.Struct(t)
}

// GetID returns the ID of the trip
func (t *TripBase) GetID() string {
	return t.ID
}

// GetClientID returns the ClientID of the trip
func (t *TripBase) GetClientID() string {
	return t.ClientID
}

// GetCreatedAt returns the creation timestamp of the trip
func (t *TripBase) GetCreatedAt() int64 {
	return t.CreatedAt
}

// GetUpdatedAt returns the last update timestamp of the trip
func (t *TripBase) GetUpdatedAt() int64 {
	return t.UpdatedAt
}

// GetTripType returns the TripType of the trip
func (t *TripBase) GetTripType() TripType {
	return t.TripType
}

// GetDeletedAt returns the deletion timestamp of the trip
func (t *TripBase) GetDeletedAt() int64 {
	return t.DeletedAt
}

// SetID sets the ID of the trip
func (t *TripBase) SetID(id string) {
	t.ID = id
}

// SetClientID sets the ClientID of the trip
func (t *TripBase) SetClientID(clientID string) {
	t.ClientID = clientID
}

// SetCreatedAt sets the creation timestamp of the trip
func (t *TripBase) SetCreatedAt(createdAt int64) {
	t.CreatedAt = createdAt
}

// SetUpdatedAt sets the last update timestamp of the trip
func (t *TripBase) SetUpdatedAt(updatedAt int64) {
	t.UpdatedAt = updatedAt
}

// SetTripType sets the TripType of the trip
func (t *TripBase) SetTripType(tripType TripType) {
	t.TripType = tripType
}

// SetDeletedAt sets the deletion timestamp of the trip
func (t *TripBase) SetDeletedAt(deletedAt int64) {
	t.DeletedAt = deletedAt
}

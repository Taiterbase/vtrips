package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/pkg/ulid"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
	validate "github.com/go-playground/validator/v10"
)

var (
	ErrInvalidPayload = fmt.Errorf("Request payload invalid")
	ErrInvalidTripID  = fmt.Errorf("TripID is required as a string")
	ErrInvalidOrgID   = fmt.Errorf("OrgID is required as a string")
	ErrTripNotFound   = fmt.Errorf("Trip not found")
	ErrOrgNotFound    = fmt.Errorf("OrgID not found")
)

type Trip interface {
	GetID() string
	GetOrgID() string
	GetCreatedAt() int64
	GetUpdatedAt() int64
	GetTripType() TripType
	GetDeletedAt() int64

	SetID(string)
	SetOrgID(string)
	SetCreatedAt(int64)
	SetUpdatedAt(int64)
	SetTripType(TripType)
	SetDeletedAt(int64)

	Validate() error
	Tokenize() [][]byte
}

type TripBase struct {
	ID    string `json:"id"`
	OrgID string `json:"org_id" validate:"required" index:"equality"`

	HousingType    HousingType `json:"housing_type" updateable:"true" index:"equality"`
	PrivacyType    PrivacyType `json:"privacy_type" updateable:"true" index:"equality"`
	TripType       TripType    `json:"trip_type" updateable:"true" index:"equality"`
	Status         TripStatus  `json:"status" updateable:"true" index:"equality"`
	VolunteerLimit int         `json:"volunteer_limit" updateable:"true" index:"equality"`
	Name           string      `json:"name" updateable:"true"`
	Description    string      `json:"description" updateable:"true"`
	Mission        string      `json:"mission" updateable:"true"`
	Price          float64     `json:"price" updateable:"true" index:"equality"`
	Currency       string      `json:"currency" updateable:"true" index:"equality"`

	Latitude   float64 `json:"latitude" updateable:"true" index:"geoposition"`
	Longtitude float64 `json:"longtitude" updateable:"true" index:"geoposition"`

	StartDate int64 `json:"start_date" updateable:"true" index:"time"`
	EndDate   int64 `json:"end_date" updateable:"true" index:"time"`
	CreatedAt int64 `json:"created_at" validate:"required" index:"time"`
	UpdatedAt int64 `json:"updated_at" updateable:"true" index:"time"`
	DeletedAt int64 `json:"deleted_at" updateable:"true" index:"time"`
}

// NewTrip creates a new Trip struct with default values
func NewTrip() Trip {
	now := time.Now().Unix()
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

// Validate validates the TripBase struct
func (t *TripBase) Validate() error {
	validate := validate.New()
	return validate.Struct(t)
}

// GetID returns the ID of the trip
func (t *TripBase) GetID() string {
	return t.ID
}

// GetOrgID returns the OrgID of the trip
func (t *TripBase) GetOrgID() string {
	return t.OrgID
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

// SetOrgID sets the OrgID of the trip
func (t *TripBase) SetOrgID(clientID string) {
	t.OrgID = clientID
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

func GetDailyBucket(timestamp int64) int64 {
	return timestamp - (timestamp % 86400)
}

func (t *TripBase) Tokenize() [][]byte {
	tokens := [][]byte{}
	typ := reflect.TypeOf(*t)
	v := reflect.ValueOf(*t)
	for i := range typ.NumField() {
		field := typ.Field(i)
		switch field.Tag.Get("index") {
		case "time":
			value := v.Field(i).Int()
			value = GetDailyBucket(value)
			token := utils.MakeKey(field.Tag.Get("json"), fmt.Sprintf("%v", value))
			tokens = append(tokens, token)
		case "geoposition":
		// todo(t8): implement
		case "equality":
			value := v.Field(i).Interface()
			if value == nil || value == "" {
				continue
			}
			token := utils.MakeKey(field.Tag.Get("json"), fmt.Sprintf("%v", value))
			tokens = append(tokens, token)
		}
	}
	return tokens
}

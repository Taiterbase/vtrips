package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Taiterbase/vtrips/apps/users/pkg/ulid"
	"github.com/Taiterbase/vtrips/apps/users/pkg/utils"
	validate "github.com/go-playground/validator/v10"
)

var (
	ErrInvalidPayload = fmt.Errorf("Request payload invalid")
	ErrInvalidUserID  = fmt.Errorf("UserID is required as a string")
	ErrInvalidOrgID   = fmt.Errorf("OrgID is required as a string")
	ErrUserNotFound   = fmt.Errorf("User not found")
	ErrOrgNotFound    = fmt.Errorf("OrgID not found")
)

type User interface {
	GetID() string
	GetOrgID() string
	GetCreatedAt() int64
	GetUpdatedAt() int64
	GetDeletedAt() int64

	SetID(string)
	SetOrgID(string)
	SetCreatedAt(int64)
	SetUpdatedAt(int64)
	SetDeletedAt(int64)

	Validate() error
	Tokenize() [][]byte
}

type UserBase struct {
	ID    string `json:"id"`
	OrgID string `json:"org_id" validate:"required" index:"equality"`

	VolunteerLimit int     `json:"volunteer_limit" updateable:"true" index:"equality"`
	Name           string  `json:"name" updateable:"true"`
	Description    string  `json:"description" updateable:"true"`
	Mission        string  `json:"mission" updateable:"true"`
	Price          float64 `json:"price" updateable:"true" index:"equality"`
	Currency       string  `json:"currency" updateable:"true" index:"equality"`

	Latitude   float64 `json:"latitude" updateable:"true" index:"geoposition"`
	Longtitude float64 `json:"longtitude" updateable:"true" index:"geoposition"`

	StartDate int64 `json:"start_date" updateable:"true" index:"time"`
	EndDate   int64 `json:"end_date" updateable:"true" index:"time"`
	CreatedAt int64 `json:"created_at" validate:"required" index:"time"`
	UpdatedAt int64 `json:"updated_at" updateable:"true" index:"time"`
	DeletedAt int64 `json:"deleted_at" updateable:"true" index:"time"`
}

// NewUser creates a new User struct with default values
func NewUser() User {
	now := time.Now().Unix()
	return &UserBase{
		ID:        ulid.Make().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the UserBase struct
func (t *UserBase) Validate() error {
	validate := validate.New()
	return validate.Struct(t)
}

// GetID returns the ID of the trip
func (t *UserBase) GetID() string {
	return t.ID
}

// GetOrgID returns the OrgID of the trip
func (t *UserBase) GetOrgID() string {
	return t.OrgID
}

// GetCreatedAt returns the creation timestamp of the trip
func (t *UserBase) GetCreatedAt() int64 {
	return t.CreatedAt
}

// GetUpdatedAt returns the last update timestamp of the trip
func (t *UserBase) GetUpdatedAt() int64 {
	return t.UpdatedAt
}

// GetDeletedAt returns the deletion timestamp of the trip
func (t *UserBase) GetDeletedAt() int64 {
	return t.DeletedAt
}

// SetID sets the ID of the trip
func (t *UserBase) SetID(id string) {
	t.ID = id
}

// SetOrgID sets the OrgID of the trip
func (t *UserBase) SetOrgID(clientID string) {
	t.OrgID = clientID
}

// SetCreatedAt sets the creation timestamp of the trip
func (t *UserBase) SetCreatedAt(createdAt int64) {
	t.CreatedAt = createdAt
}

// SetUpdatedAt sets the last update timestamp of the trip
func (t *UserBase) SetUpdatedAt(updatedAt int64) {
	t.UpdatedAt = updatedAt
}

// SetDeletedAt sets the deletion timestamp of the trip
func (t *UserBase) SetDeletedAt(deletedAt int64) {
	t.DeletedAt = deletedAt
}

func GetDailyBucket(timestamp int64) int64 {
	return timestamp - (timestamp % 86400)
}

func (t *UserBase) Tokenize() [][]byte {
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

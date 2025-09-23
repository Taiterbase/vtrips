package models

import (
	"fmt"
	"reflect"
	"time"

	validate "github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidPayload = fmt.Errorf("Request payload invalid")
	ErrInvalidUserID  = fmt.Errorf("UserID is required as a string")
	ErrInvalidOrgID   = fmt.Errorf("OrgID is required as a string")
	ErrUserNotFound   = fmt.Errorf("User not found")
	ErrOrgNotFound    = fmt.Errorf("OrgID not found")
)

type User struct {
	ID            string `json:"id"`
    Username      string `json:"username" db:"username" index:"equality"`
    Hash          string `json:"hash" db:"hash"`
    Contact       string `json:"contact" db:"contact" index:"equality"`
    ContactMethod string `json:"contact_method" db:"contact_method"`
	DOB           string `json:"dob" db:"dob"`

	CreatedAt int64 `json:"created_at" validate:"required" index:"time"`
	UpdatedAt int64 `json:"updated_at" updateable:"true" index:"time"`
	DeletedAt int64 `json:"deleted_at" updateable:"true" index:"time"`
}

// NewUser creates a new User struct with default values
func NewUser() *User {
	now := time.Now().Unix()
	return &User{
		ID:        ulid.Make().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the User struct
func (t *User) Validate() error {
	validate := validate.New()
	return validate.Struct(t)
}

// Accessors to match storage expectations (mirror trips model style)
func (t *User) GetID() string { return t.ID }
func (t *User) GetUpdatedAt() int64 { return t.UpdatedAt }
func (t *User) SetUpdatedAt(ts int64) { t.UpdatedAt = ts }
func (t *User) GetUsername() string { return t.Username }
func (t *User) GetHash() string { return t.Hash }

func GetDailyBucket(timestamp int64) int64 {
	return timestamp - (timestamp % 86400)
}

func (t *User) Tokenize() [][]byte {
	tokens := [][]byte{}
	typ := reflect.TypeOf(*t)
	v := reflect.ValueOf(*t)
	for i := range typ.NumField() {
		field := typ.Field(i)
		switch field.Tag.Get("index") {
		case "time":
			value := v.Field(i).Int()
			value = GetDailyBucket(value)
			token := MakeKey(field.Tag.Get("json"), fmt.Sprintf("%v", value))
			tokens = append(tokens, token)
		case "geoposition":
		// todo(t8): implement
		case "equality":
			value := v.Field(i).Interface()
			if value == nil || value == "" {
				continue
			}
			token := MakeKey(field.Tag.Get("json"), fmt.Sprintf("%v", value))
			tokens = append(tokens, token)
		}
	}
	return tokens
}

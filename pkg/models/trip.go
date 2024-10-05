package models

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TripType int

const (
	OtherTrip TripType = iota
	LocalTrip
	DomesticTrip
	InternationalTrip
)

func (t TripType) String() string {
	return [...]string{"local", "domestic", "international"}[t]
}

func (t TripType) Int() int {
	return int(t)
}

func (t TripType) MarshalJSON() ([]byte, error) {
	switch t {
	case LocalTrip:
		return []byte("local"), nil
	case DomesticTrip:
		return []byte("domestic"), nil
	case InternationalTrip:
		return []byte("international"), nil
	default:
		return []byte("other"), nil
	}
}

func (t TripType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "local":
		t = LocalTrip
	case "domestic":
		t = DomesticTrip
	case "international":
		t = InternationalTrip
	default:
		t = LocalTrip
	}
	return nil
}

func (t TripType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: t.String()}, nil
}

func (t *TripType) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	switch av := av.(type) {
	case *types.AttributeValueMemberS:
		switch av.Value {
		case "local":
			*t = LocalTrip
		case "domestic":
			*t = DomesticTrip
		case "international":
			*t = InternationalTrip
		default:
			*t = OtherTrip
		}
	default:
		*t = TripType(OtherTrip)
	}
	return nil
}

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
}

type TripBase struct {
	PK        string `json:"-" dynamodbav:"pk"`
	SK        string `json:"-" dynamodbav:"sk"`
	ID        string `json:"id" dynamodbav:"id"`
	ClientID  string `json:"client_id" dynamodbav:"client_id"`
	CreatedAt int64  `json:"created_at" dynamodbav:",unixtime"`
	UpdatedAt int64  `json:"updated_at" dynamodbav:",unixtime"`
	DeletedAt int64  `json:"deleted_at" dynamodbav:",unixtime"`

	HousingType HousingType `json:"housing_type" dynamodbav:"housing_type"`
	PrivacyType PrivacyType `json:"privacy_type" dynamodbav:"privacy_type"`
	TripType    TripType    `json:"trip_type" dynamodbav:"trip_type"`
}

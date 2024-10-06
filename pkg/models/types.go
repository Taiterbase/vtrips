package models

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TripType int

const (
	OtherTrip TripType = iota
	LocalTrip
	DomesticTrip
	InternationalTrip
)

func (t TripType) MarshalJSON() ([]byte, error) {
	var val string
	switch t {
	case LocalTrip:
		val = "local"
	case DomesticTrip:
		val = "domestic"
	case InternationalTrip:
		val = "international"
	default:
		val = "other"
	}
	return json.Marshal(val)
}

func (t *TripType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch val {
	case "local":
		*t = LocalTrip
	case "domestic":
		*t = DomesticTrip
	case "international":
		*t = InternationalTrip
	default:
		*t = OtherTrip
	}

	return nil
}

func (t TripType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var val string
	switch t {
	case LocalTrip:
		val = "local"
	case DomesticTrip:
		val = "domestic"
	case InternationalTrip:
		val = "international"
	default:
		val = "other"
	}
	return &types.AttributeValueMemberS{Value: val}, nil
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

type PrivacyType int

const (
	OtherPrivacy PrivacyType = iota
	SharedPrivacy
	PrivatePrivacy
	CompletePrivacy
)

func (p PrivacyType) MarshalJSON() ([]byte, error) {
	var val string
	switch p {
	case SharedPrivacy:
		val = "shared"
	case PrivatePrivacy:
		val = "private"
	case CompletePrivacy:
		val = "complete"
	default:
		val = "other"
	}
	return json.Marshal(val)
}

func (p *PrivacyType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch val {
	case "shared":
		*p = SharedPrivacy
	case "private":
		*p = PrivatePrivacy
	case "complete":
		*p = CompletePrivacy
	default:
		*p = OtherPrivacy
	}

	return nil
}

func (p PrivacyType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var val string
	switch p {
	case SharedPrivacy:
		val = "shared"
	case PrivatePrivacy:
		val = "private"
	case CompletePrivacy:
		val = "complete"
	case OtherPrivacy:
		val = "other"
	default:
		val = "other"
	}
	return &types.AttributeValueMemberS{Value: val}, nil
}

func (p *PrivacyType) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	switch av := av.(type) {
	case *types.AttributeValueMemberS:
		switch av.Value {
		case "shared":
			*p = SharedPrivacy
		case "private":
			*p = PrivatePrivacy
		case "complete":
			*p = CompletePrivacy
		case "other":
			*p = OtherPrivacy
		default:
			*p = OtherPrivacy
		}
	default:
		*p = PrivacyType(OtherPrivacy)
	}
	return nil
}

type HousingType int

const (
	OtherHousing HousingType = iota
	CampingHousing
	HostelHousing
	HotelHousing
	DormitoryHousing
	ApartmentHousing
	HouseHousing
)

func (h HousingType) MarshalJSON() ([]byte, error) {
	var val string
	switch h {
	case CampingHousing:
		val = "camping"
	case HostelHousing:
		val = "hostel"
	case HotelHousing:
		val = "hotel"
	case DormitoryHousing:
		val = "dormitory"
	case ApartmentHousing:
		val = "apartment"
	case HouseHousing:
		val = "house"
	default:
		val = "other"
	}
	return json.Marshal(val)
}

func (h *HousingType) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch val {
	case "camping":
		*h = CampingHousing
	case "hostel":
		*h = HostelHousing
	case "hotel":
		*h = HotelHousing
	case "dormitory":
		*h = DormitoryHousing
	case "apartment":
		*h = ApartmentHousing
	case "house":
		*h = HouseHousing
	default:
		*h = OtherHousing
	}

	return nil
}

func (h HousingType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	var val string
	switch h {
	case CampingHousing:
		val = "camping"
	case HostelHousing:
		val = "hostel"
	case HotelHousing:
		val = "hotel"
	case DormitoryHousing:
		val = "dormitory"
	case ApartmentHousing:
		val = "apartment"
	case HouseHousing:
		val = "house"
	case OtherHousing:
		val = "other"
	default:
		val = "other"
	}

	return &types.AttributeValueMemberS{Value: val}, nil
}

func (h *HousingType) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	switch av := av.(type) {
	case *types.AttributeValueMemberS:
		switch av.Value {
		case "camping":
			*h = CampingHousing
		case "hostel":
			*h = HostelHousing
		case "hotel":
			*h = HotelHousing
		case "dormitory":
			*h = DormitoryHousing
		case "apartment":
			*h = ApartmentHousing
		case "house":
			*h = HouseHousing
		case "other":
			*h = OtherHousing
		default:
			*h = OtherHousing
		}
	default:
		*h = HousingType(OtherHousing)
	}
	return nil
}

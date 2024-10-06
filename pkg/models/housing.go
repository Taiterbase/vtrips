package models

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

func (h HousingType) String() string {
	return [...]string{"camping", "hostel", "hotel", "dormitory", "apartment", "house", "other"}[h]
}

func (h HousingType) Int() int {
	return int(h)
}

func (h HousingType) MarshalJSON() ([]byte, error) {
	switch h {
	case CampingHousing:
		return json.Marshal("camping")
	case HostelHousing:
		return json.Marshal("hostel")
	case HotelHousing:
		return json.Marshal("hotel")
	case DormitoryHousing:
		return json.Marshal("dormitory")
	case ApartmentHousing:
		return json.Marshal("apartment")
	case HouseHousing:
		return json.Marshal("house")
	case OtherHousing:
		return json.Marshal("other")
	default:
		return json.Marshal("other")
	}
}

func (h HousingType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "camping":
		h = CampingHousing
	case "hostel":
		h = HostelHousing
	case "hotel":
		h = HotelHousing
	case "dormitory":
		h = DormitoryHousing
	case "apartment":
		h = ApartmentHousing
	case "house":
		h = HouseHousing
	case "other":
		h = OtherHousing
	default:
		h = OtherHousing
	}
	return nil
}

func (h HousingType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: h.String()}, nil
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

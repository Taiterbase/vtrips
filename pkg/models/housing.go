package models

import (
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

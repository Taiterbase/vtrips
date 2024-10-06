package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

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
	return [...]string{"camping", "hostel", "hotel", "dormitory", "apartment", "house"}[h]
}

func (h HousingType) Int() int {
	return int(h)
}

func (h HousingType) MarshalJSON() ([]byte, error) {
	switch h {
	case CampingHousing:
		return []byte("camping"), nil
	case HostelHousing:
		return []byte("hostel"), nil
	case HotelHousing:
		return []byte("hotel"), nil
	case DormitoryHousing:
		return []byte("dormitory"), nil
	case ApartmentHousing:
		return []byte("apartment"), nil
	case HouseHousing:
		return []byte("house"), nil
	default:
		return []byte("other"), nil
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
		default:
			*h = OtherHousing
		}
	default:
		*h = HousingType(OtherHousing)
	}
	return nil
}

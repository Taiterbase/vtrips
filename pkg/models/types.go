package models

type TripType string

const (
	OtherTrip         TripType = "other"
	LocalTrip         TripType = "local"
	DomesticTrip      TripType = "domestic"
	InternationalTrip TripType = "international"
)

type PrivacyType string

const (
	OtherPrivacy    PrivacyType = "other"
	SharedPrivacy   PrivacyType = "shared"
	PrivatePrivacy  PrivacyType = "private"
	CompletePrivacy PrivacyType = "complete"
)

type HousingType string

const (
	OtherHousing     HousingType = "other"
	CampingHousing   HousingType = "camping"
	HostelHousing    HousingType = "hostel"
	HotelHousing     HousingType = "hotel"
	DormitoryHousing HousingType = "dormitory"
	ApartmentHousing HousingType = "apartment"
	HouseHousing     HousingType = "house"
)

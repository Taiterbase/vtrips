package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type PrivacyType int

const (
	OtherPrivacy PrivacyType = iota
	SharedPrivacy
	PrivatePrivacy
	CompletePrivacy
)

func (p PrivacyType) String() string {
	return [...]string{"shared", "private", "complete", "other"}[p]
}

func (p PrivacyType) Int() int {
	return int(p)
}

func (p PrivacyType) MarshalJSON() ([]byte, error) {
	switch p {
	case SharedPrivacy:
		return []byte("shared"), nil
	case PrivatePrivacy:
		return []byte("private"), nil
	case CompletePrivacy:
		return []byte("complete"), nil
	case OtherPrivacy:
		return []byte("other"), nil
	default:
		return []byte("other"), nil
	}
}

func (p PrivacyType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "shared":
		p = SharedPrivacy
	case "private":
		p = PrivatePrivacy
	case "complete":
		p = CompletePrivacy
	default:
		p = SharedPrivacy
	}
	return nil
}

func (p PrivacyType) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{Value: p.String()}, nil
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

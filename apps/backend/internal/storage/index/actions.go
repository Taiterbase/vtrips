package index

import (
	"github.com/Taiterbase/vtrips/pkg/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func GetTripWriteActions() []WriteAction {
	return []WriteAction{
		{
			Add: func(tx *dynamodb.TransactWriteItemsInput, trip models.Trip) error {
				return nil
			},
			Update: func(tx *dynamodb.TransactWriteItemsInput, trip models.Trip, updatedFields []string) error {
				return nil
			},
			Delete: func(tx *dynamodb.TransactWriteItemsInput, trip models.Trip) error {
				return nil
			},
		},
	}
}

package index

import (
	"github.com/Taiterbase/vtrips/pkg/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// WriteAction is a set of functions that define transact write behavior when managing Trips
type WriteAction struct {
	Add    func(*dynamodb.TransactWriteItemsInput, string, models.Trip) error
	Update func(*dynamodb.TransactWriteItemsInput, string, models.Trip, []string) error
	Delete func(*dynamodb.TransactWriteItemsInput, string, models.Trip) error
}

// CascadeAction is a set of functions that define behavior around operations
// in the system that should cascade.
type CascadeAction struct {
	Add    func([]models.Trip, models.Trip) ([]*dynamodb.TransactWriteItemsInput, error)
	Update func([]models.Trip, models.Trip, []string) ([]*dynamodb.TransactWriteItemsInput, error)
	Delete func([]models.Trip, models.Trip) ([]*dynamodb.TransactWriteItemsInput, error)
}

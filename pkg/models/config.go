package models

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// WriteAction is a set of functions that define transact write behavior when managing Trips
type WriteAction struct {
	Add    func(*dynamodb.TransactWriteItemsInput, string, Trip) error
	Update func(*dynamodb.TransactWriteItemsInput, string, Trip, Trip) error
	Delete func(*dynamodb.TransactWriteItemsInput, string, Trip) error
}

// CascadeAction is a set of functions that define behavior around operations
// in the system that should cascade.
type CascadeAction struct {
	Add    func([]Trip, Trip) ([]*dynamodb.TransactWriteItemsInput, error)
	Update func([]Trip, Trip, []string) ([]*dynamodb.TransactWriteItemsInput, error)
	Delete func([]Trip, Trip) ([]*dynamodb.TransactWriteItemsInput, error)
}

// IndexInterface is an interface that defines getters for actions to take when modifying a trip
type Indexer interface {
	GetCascadeActions() []CascadeAction
	GetWriteActions() []WriteAction
}

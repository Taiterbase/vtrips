package models

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// WriteAction is a set of functions that define transact write behavior when managing Trips
type WriteAction[T Trip] struct {
	Add    func(*dynamodb.TransactWriteItemsInput, string, T) error
	Update func(*dynamodb.TransactWriteItemsInput, string, T, T) error
	Delete func(*dynamodb.TransactWriteItemsInput, string, T) error
}

// CascadeAction is a set of functions that define behavior around operations
// in the system that should cascade.
type CascadeAction[T Trip] struct {
	Add    func([]T, string, T) ([]*dynamodb.TransactWriteItemsInput, error)
	Update func([]T, string, T, []string) ([]*dynamodb.TransactWriteItemsInput, error)
	Delete func([]T, string, T) ([]*dynamodb.TransactWriteItemsInput, error)
}

// IndexInterface is an interface that defines getters for actions to take when modifying a trip
type IndexInterface[T Trip] interface {
	GetCascadeActions() []CascadeAction[T]
	GetWriteActions() WriteAction[T]
}

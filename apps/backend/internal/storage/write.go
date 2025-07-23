package storage

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/pkg/models"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

// writeTokens writes the trip to the inverted index for each of its fields
func writeTokens(c echo.Context, batch *pebble.Batch, trip models.Trip) (err error) {
	startTime := time.Now()
	defer func(err error) {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		if err != nil {
			c.Logger().Debugj(log.JSON{"WriteTokens": duration, "error": err.Error()})
		}
	}(err)

	tripID := trip.GetID()
	for _, key := range trip.Tokenize() {
		existingData, closer, err := batch.Get(key)
		if err != nil && err != pebble.ErrNotFound {
			return err
		}
		var existingIDs []string
		if err == nil {
			defer closer.Close()
			err = json.Unmarshal(existingData, &existingIDs)
			if err != nil {
				return err
			}
		}

		if !slices.Contains(existingIDs, tripID) {
			existingIDs = append(existingIDs, tripID)
		}

		updatedData, err := json.Marshal(existingIDs)
		if err != nil {
			return err
		}

		err = batch.Set(key, updatedData, pebble.Sync)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTrip creates a new trip in the database and writes to an inverted index for each of its fields
func CreateTrip(c echo.Context, trip models.Trip) error {
	batch := Client.NewIndexedBatch()
	key := utils.MakeKey("trip_id", trip.GetID())
	tripBytes, err := json.Marshal(trip)
	if err != nil {
		return err
	}
	err = batch.Set(key, tripBytes, pebble.Sync)
	if err != nil {
		return err
	}
	err = writeTokens(c, batch, trip)
	if err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// DeleteTrip deletes a trip from the database and removes it from the inverted index
func DeleteTrip(c echo.Context, trip models.Trip) (err error) {
	batch := Client.NewIndexedBatch()
	tripID := trip.GetID()
	key := utils.MakeKey("trip_id", tripID)

	var existingTrip models.Trip
	tripBytes, closer, err := batch.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil
		}
		return err
	}
	defer closer.Close()

	err = json.Unmarshal(tripBytes, &existingTrip)
	if err != nil {
		return err
	}

	existingTokens := existingTrip.Tokenize()
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		if err != nil {
			c.Logger().Debugj(log.JSON{"delete_trip": duration, "error": err.Error()})
		}
	}()

	for _, tokenKey := range existingTokens {
		existingData, closer, err := batch.Get(tokenKey)
		if err != nil {
			if err != pebble.ErrNotFound {
				return err
			}
			continue
		}
		defer closer.Close()

		var existingIDs []string
		err = json.Unmarshal(existingData, &existingIDs)
		if err != nil {
			return err
		}

		existingIDs = utils.Remove(existingIDs, tripID)
		if len(existingIDs) == 0 {
			err = batch.Delete(tokenKey, pebble.Sync)
		} else {
			updatedData, err := json.Marshal(existingIDs)
			if err != nil {
				return err
			}
			err = batch.Set(tokenKey, updatedData, pebble.Sync)
		}
		if err != nil {
			return err
		}
	}

	err = batch.Delete(key, pebble.Sync)
	if err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// UpdateTrip updates a trip in the database and updates the inverted index for each of its fields
func UpdateTrip(c echo.Context, trip models.Trip) (err error) {
	batch := Client.NewIndexedBatch()
	tripID := trip.GetID()
	key := utils.MakeKey("trip_id", tripID)
	var existingTrip models.Trip
	tripBytes, closer, err := batch.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return models.ErrTripNotFound
		}
		return err
	}
	defer closer.Close()
	err = json.Unmarshal(tripBytes, &existingTrip)
	if err != nil {
		return err
	}
	existingTokens := existingTrip.Tokenize()
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		if err != nil {
			c.Logger().Debugj(log.JSON{"update_trip": duration, "error": err.Error()})
		}
	}()

	// naive update
	for _, tokenKey := range existingTokens {
		existingData, closer, err := batch.Get(tokenKey)
		if err != nil {
			if err != pebble.ErrNotFound {
				return err
			}
			continue
		}
		defer closer.Close()
		var existingIDs []string
		err = json.Unmarshal(existingData, &existingIDs)
		if err != nil {
			return err
		}
		existingIDs = utils.Remove(existingIDs, tripID)
		if len(existingIDs) == 0 {
			err = batch.Delete(tokenKey, pebble.Sync)
		} else {
			updatedData, err := json.Marshal(existingIDs)
			if err != nil {
				return err
			}
			err = batch.Set(tokenKey, updatedData, pebble.Sync)
		}
		if err != nil {
			return err
		}
	}
	err = writeTokens(c, batch, trip)
	if err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

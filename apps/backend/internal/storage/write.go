package storage

import (
	"encoding/json"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/pkg/models"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
)

// CreateTrip creates a new trip in the database and writes to an inverted index for each of its fields
func CreateTrip(c echo.Context, trip models.Trip) error {
	batch := Client.NewBatch()
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

func writeTokens(c echo.Context, batch *pebble.Batch, trip models.Trip) (err error) {
	startTime := time.Now()
	defer func(err error) {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		if err != nil {
			c.Logger().Debugf("writeTokens: %v", err)
		}
		c.Logger().Debugf("writeTokens: %v", duration)
	}(err)

	tripID := trip.GetID()
	for _, key := range trip.Tokenize() {
		existingData, closer, err := Client.Get(key)
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

		if !utils.Contains(existingIDs, tripID) {
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

package storage

import (
	"encoding/json"

	"github.com/Taiterbase/vtrips/apps/backend/pkg/models"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
)

func ReadTrip(c echo.Context, tripID string) (models.Trip, error) {
	var trip models.TripBase
	c.Logger().Debugf("Reading trip: %v", tripID)
	tripKey := utils.MakeKey("trip_id", tripID)
	tripBytes, closer, err := Client.Get(tripKey)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, models.ErrTripNotFound
		}
		return nil, err
	}
	defer closer.Close()
	err = json.Unmarshal(tripBytes, &trip)
	return &trip, err
}

func ReadTrips(c echo.Context, tripID []string) ([]models.Trip, error) {
	var trips []models.Trip
	for _, tripID := range tripID {
		trip, err := ReadTrip(c, tripID)
		if err != nil {
			return nil, err
		}
		if trip == nil {
			continue
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

func GetIDsFromToken(c echo.Context, key []byte) ([]string, error) {
	data, closer, err := Client.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()
	var ids []string
	err = json.Unmarshal(data, &ids)
	return ids, err
}

package storage

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/trips/pkg/models"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
)

// BitmapForToken returns the decoded bitmap for tokenKey or an empty bitmap.
func BitmapForToken(tokenKey []byte) (*roaring64.Bitmap, error) {
	v, closer, err := Client.Get(tokenKey)
	if err == pebble.ErrNotFound {
		return roaring64.New(), nil
	}
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return decode(v)
}

func ReadTrip(c echo.Context, tripID string) (models.Trip, error) {
	var trip models.TripBase
	tripKey := models.MakeKey("trip_id", tripID)
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

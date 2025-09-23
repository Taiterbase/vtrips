package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/trips/internal/storage"
	"github.com/Taiterbase/vtrips/apps/trips/pkg/models"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func DatabaseDebug(c echo.Context) error {
	iter, err := storage.Client.NewIter(&pebble.IterOptions{})
	if err != nil {
		return err
	}
	defer iter.Close()

	dbContents := make(map[string]any)
	valid := iter.First()
	for valid {
		key := string(iter.Key())
		value := iter.Value()
		var prettyValue any
		if err := json.Unmarshal(value, &prettyValue); err != nil {
			prettyValue = string(value)
		}
		dbContents[key] = prettyValue
		valid = iter.Next()
	}

	return c.JSON(http.StatusOK, map[string]any{
		"database_contents": dbContents,
	})
}

func CreateTrip(c echo.Context) error {
	trip := models.NewTrip()
	err := c.Bind(&trip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = trip.Validate(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = storage.CreateTrip(c, trip)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trip)
}

func getTrip(c echo.Context, orgID, tripID string) (models.Trip, error) {
	numID, ok, err := storage.Lookup(storage.Client, tripID)
	if err != nil {
		return nil, err
	}
	if !ok { // unknown ULID
		return nil, models.ErrTripNotFound
	}

	orgToken := models.MakeKey("org_id", orgID)
	bm, err := storage.BitmapForToken(orgToken)
	if err != nil {
		return nil, err
	}
	if !bm.Contains(numID) {
		return nil, models.ErrTripNotFound
	}

	return storage.ReadTrip(c, tripID)
}

func GetTrip(c echo.Context) error {
	orgID := c.QueryParam("org_id")
	tripID := c.Param("trip_id")
	if orgID == "" || tripID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrInvalidTripID)
	}
	trip, err := getTrip(c, orgID, tripID)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, trip)
	case models.ErrTripNotFound:
		return c.JSON(http.StatusNotFound, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func GetTrips(c echo.Context) error {
	scannedCount := 0
	bitmapMap := make(map[string]*roaring64.Bitmap)
	for key, vals := range c.QueryParams() {
		for _, v := range vals {
			tk := models.MakeKey(key, v)
			bm, err := storage.BitmapForToken(tk)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			if orBm, ok := bitmapMap[key]; !ok {
				bitmapMap[key] = bm
			} else {
				// if we've seen this key before in the query parameters, take a running union of the bm
				bitmapMap[key] = roaring64.Or(bm, orBm)
			}
		}
	}

	var bms []*roaring64.Bitmap
	for _, bm := range bitmapMap {
		scannedCount += int(bm.GetCardinality())
		bms = append(bms, bm)
	}
	intersection := roaring64.FastAnd(bms...)
	if intersection.IsEmpty() {
		return c.JSON(http.StatusOK, []models.Trip{})
	}
	c.Logger().Debugj(log.JSON{"message": "we have an intersection", "intersection": intersection})

	var trips []models.Trip
	it := intersection.Iterator()
	for it.HasNext() {
		numID := it.Next()
		ulid, ok, err := storage.Reverse(storage.Client, numID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if !ok {
			continue
		} // should not happen
		t, err := storage.ReadTrip(c, ulid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		trips = append(trips, t)
	}
	return c.JSON(http.StatusOK, log.JSON{
		"trips":         trips,
		"count":         len(trips),
		"scanned_count": scannedCount,
	})
}

func UpdateTrip(c echo.Context) error {
	var (
		tripID = c.Param("trip_id")
		orgID  = c.QueryParam("org_id")
	)
	trip, err := getTrip(c, orgID, tripID)
	switch err {
	case models.ErrTripNotFound:
		return c.JSON(http.StatusNotFound, err.Error())
	case nil:
		break
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err = c.Bind(&trip); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	trip.SetUpdatedAt(time.Now().Unix())
	err = storage.UpdateTrip(c, trip)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func DeleteTrip(c echo.Context) error {
	var (
		orgID  = c.QueryParam("org_id")
		tripID = c.Param("trip_id")
	)
	trip, err := getTrip(c, orgID, tripID)
	switch err {
	case models.ErrTripNotFound:
		return c.JSON(http.StatusNotFound, err.Error())
	case nil:
		break
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = storage.DeleteTrip(c, trip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tripID)
}

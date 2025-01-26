package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Taiterbase/vtrips/apps/backend/internal/storage"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/models"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
)

func DatabaseDebug(c echo.Context) error {
	iter, err := storage.Client.NewIter(&pebble.IterOptions{})
	if err != nil {
		return err
	}
	defer iter.Close()

	dbContents := make(map[string]interface{})
	valid := iter.First()
	for valid {
		key := string(iter.Key())
		value := iter.Value()
		var prettyValue interface{}
		if err := json.Unmarshal(value, &prettyValue); err != nil {
			prettyValue = string(value)
		}
		dbContents[key] = prettyValue
		valid = iter.Next()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"database_contents": dbContents,
	})
}

func CreateTrip(c echo.Context) error {
	trip := models.NewTrip()
	err := c.Bind(&trip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
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

func GetTrip(c echo.Context) error {
	var (
		org    = c.QueryParam("org_id")
		tripID = c.Param("trip_id")
	)
	if tripID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrInvalidTripID)
	}
	if org == "" {
		return c.JSON(http.StatusBadRequest, models.ErrInvalidOrgID)
	}

	trip, err := getTrip(c, org, tripID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, trip)
}

func getTrip(c echo.Context, org, tripID string) (models.Trip, error) {
	orgs, err := storage.GetIDsFromToken(c, []byte(fmt.Sprintf("org_id:%s", org)))
	if err != nil {
		return nil, err
	}
	// take the intersection of the two lists
	trips := utils.Intersection(orgs, []string{tripID})
	if len(trips) == 0 {
		return nil, models.ErrTripNotFound
	}
	return storage.ReadTrip(c, trips[0])
}

func UpdateTrip(c echo.Context) error {
	var (
		tripID = c.Param("trip_id")
		org    = c.QueryParam("org_id")
	)
	trip, err := getTrip(c, org, tripID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err = c.Bind(&trip); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	trip.SetUpdatedAt(time.Now().Unix())
	// write/apply the update to the tokens, handling removal and addition from lists

	return c.JSON(http.StatusOK, trip)
}

func DeleteTrip(c echo.Context) error {
	var (
		_      = c.QueryParam("org_id")
		tripID = c.Param("trip_id")
	)

	return c.JSON(http.StatusOK, tripID)
}

func GetTrips(c echo.Context) error {
	params := c.QueryParams()
	tokens := [][]byte{}
	for key, values := range params {
		if len(values) == 0 {
			continue
		}
		for _, value := range values {
			// if there are multiple values, this implies an OR condition, which means a set union and not a set intersection
			tokens = append(tokens, []byte(fmt.Sprintf("%s:%s", key, value)))
		}
	}

	// naive set intersection implementation
	tokenList := map[string][]string{}
	for _, token := range tokens {
		tokenIDs, err := storage.GetIDsFromToken(c, token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		tokenList[string(token)] = tokenIDs
	}
	intersection := []string{}
	for _, ids := range tokenList {
		if len(intersection) == 0 {
			intersection = ids
			continue
		}
		intersection = utils.Intersection(intersection, ids)
		if len(intersection) == 0 {
			return c.JSON(http.StatusOK, []models.Trip{})
		}
	}
	trips, err := storage.ReadTrips(c, intersection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, trips)
}

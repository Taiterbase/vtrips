package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/users/internal/storage"
	"github.com/Taiterbase/vtrips/apps/users/pkg/models"
	"github.com/Taiterbase/vtrips/apps/users/pkg/utils"
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

func CreateUser(c echo.Context) error {
	trip := models.NewUser()
	err := c.Bind(&trip)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err = trip.Validate(); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = storage.CreateUser(c, trip)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, trip)
}

func getUser(c echo.Context, orgID, tripID string) (models.User, error) {
	numID, ok, err := storage.Lookup(storage.Client, tripID)
	if err != nil {
		return nil, err
	}
	if !ok { // unknown ULID
		return nil, models.ErrUserNotFound
	}

	orgToken := utils.MakeKey("org_id", orgID)
	bm, err := storage.BitmapForToken(orgToken)
	if err != nil {
		return nil, err
	}
	if !bm.Contains(numID) {
		return nil, models.ErrUserNotFound
	}

	return storage.ReadUser(c, tripID)
}

func GetUser(c echo.Context) error {
	orgID := c.QueryParam("org_id")
	tripID := c.Param("trip_id")
	if orgID == "" || tripID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrInvalidUserID)
	}
	trip, err := getUser(c, orgID, tripID)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, trip)
	case models.ErrUserNotFound:
		return c.JSON(http.StatusNotFound, err.Error())
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func GetUsers(c echo.Context) error {
	scannedCount := 0
	bitmapMap := make(map[string]*roaring64.Bitmap)
	for key, vals := range c.QueryParams() {
		for _, v := range vals {
			tk := utils.MakeKey(key, v)
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
		return c.JSON(http.StatusOK, []models.User{})
	}
	c.Logger().Debugj(log.JSON{"message": "we have an intersection", "intersection": intersection})

	var trips []models.User
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
		t, err := storage.ReadUser(c, ulid)
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

func UpdateUser(c echo.Context) error {
	var (
		tripID = c.Param("trip_id")
		orgID  = c.QueryParam("org_id")
	)
	trip, err := getUser(c, orgID, tripID)
	switch err {
	case models.ErrUserNotFound:
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
	err = storage.UpdateUser(c, trip)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, nil)
}

func DeleteUser(c echo.Context) error {
	var (
		orgID  = c.QueryParam("org_id")
		tripID = c.Param("trip_id")
	)
	trip, err := getUser(c, orgID, tripID)
	switch err {
	case models.ErrUserNotFound:
		return c.JSON(http.StatusNotFound, err.Error())
	case nil:
		break
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = storage.DeleteUser(c, trip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tripID)
}

// Auth handlers (stubs)
func LoginHandler(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]any{"message": "login not implemented"})
}

func SignUpHandler(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]any{"message": "signup not implemented"})
}

func LogoutHandler(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]any{"message": "logout not implemented"})
}

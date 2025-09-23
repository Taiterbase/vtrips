package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/users/internal/auth"
	"github.com/Taiterbase/vtrips/apps/users/internal/storage"
	"github.com/Taiterbase/vtrips/apps/users/pkg/models"
	"github.com/cockroachdb/pebble"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
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

func getUser(c echo.Context, orgID, tripID string) (*models.User, error) {
	numID, ok, err := storage.Lookup(storage.Client, tripID)
	if err != nil {
		return nil, err
	}
	if !ok { // unknown ULID
		return nil, models.ErrUserNotFound
	}

	orgToken := models.MakeKey("org_id", orgID)
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
		return c.JSON(http.StatusOK, []models.User{})
	}
	c.Logger().Debugj(log.JSON{"message": "we have an intersection", "intersection": intersection})

	var trips []*models.User
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

func SignUpHandler(c echo.Context) error {
	username := c.FormValue("username")
	contact := c.FormValue("contact")
	password := c.FormValue("password")
	if username == "" || contact == "" || password == "" {
		return c.JSON(http.StatusBadRequest, "missing required fields")
	}

	// Enforce unique username via equality index bitmap
	unameToken := models.MakeKey("username", username)
	existingBm, err := storage.BitmapForToken(unameToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if !existingBm.IsEmpty() {
		return c.JSON(http.StatusConflict, "username already exists")
	}

	// hash password
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	u := models.NewUser()
	u.Username = username
	u.Contact = contact
	u.ContactMethod = "email"
	u.Hash = string(hashBytes)

	if err := u.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := storage.CreateUser(c, u); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_ = auth.SetCookieWithJWT(c.Response().Writer, u.ID, u.Username, "")
	return c.JSON(http.StatusOK, echo.Map{"id": u.ID, "username": u.Username})
}

func LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "" || password == "" {
		return c.JSON(http.StatusBadRequest, "missing credentials")
	}

	// Lookup by username index
	token := models.MakeKey("username", username)
	bm, err := storage.BitmapForToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if bm.IsEmpty() {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}
	// Use the first match (usernames should be unique)
	it := bm.Iterator()
	if !it.HasNext() {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}
	numID := it.Next()
	ulid, ok, err := storage.Reverse(storage.Client, numID)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}
	usr, err := storage.ReadUser(c, ulid)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}
	// compare hash
	if err := bcrypt.CompareHashAndPassword([]byte(usr.GetHash()), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid credentials")
	}

	_ = auth.SetCookieWithJWT(c.Response().Writer, usr.GetID(), usr.GetUsername(), "")
	return c.JSON(http.StatusOK, echo.Map{"id": usr.GetID(), "username": usr.GetUsername()})
}

func LogoutHandler(c echo.Context) error {
	// Read current token and revoke by user ID
	cookie, err := c.Cookie(auth.AuthCookieName)
	if err == nil && cookie != nil && cookie.Value != "" {
		if claims, err := auth.Validate(cookie.Value); err == nil {
			if m, ok := claims.(jwt.MapClaims); ok {
				if sub, ok := m["sub"].(string); ok && sub != "" {
					_ = storage.SetRevokedBefore(sub, time.Now().Unix())
				}
			}
		}
	}
	auth.InvalidateCookie(c.Response().Writer)
	return c.JSON(http.StatusOK, echo.Map{"ok": true})
}

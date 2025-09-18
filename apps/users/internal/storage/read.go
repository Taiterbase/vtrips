package storage

import (
	"encoding/json"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/users/pkg/models"
	"github.com/Taiterbase/vtrips/apps/users/pkg/utils"
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

func ReadUser(c echo.Context, userID string) (models.User, error) {
	var user models.UserBase
	userKey := utils.MakeKey("user_id", userID)
	userBytes, closer, err := Client.Get(userKey)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	defer closer.Close()
	err = json.Unmarshal(userBytes, &user)
	return &user, err
}

func ReadUsers(c echo.Context, userID []string) ([]models.User, error) {
	var users []models.User
	for _, userID := range userID {
		user, err := ReadUser(c, userID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
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

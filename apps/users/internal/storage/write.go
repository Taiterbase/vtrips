package storage

import (
	"encoding/json"
	"slices"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/users/pkg/models"
	"github.com/cockroachdb/pebble"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func decode(b []byte) (*roaring64.Bitmap, error) {
	rb := roaring64.New()
	if len(b) == 0 {
		return rb, nil
	}
	if err := rb.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return rb, nil
}

func encode(rb *roaring64.Bitmap) ([]byte, error) {
	rb.RunOptimize()
	return rb.MarshalBinary()
}

func writeTokens(c echo.Context, batch *pebble.Batch, user *models.User, numericID uint64) (err error) {
	defer func() {
		if err != nil {
			c.Logger().Errorj(log.JSON{"err": err})
		}
	}()

	for _, key := range user.Tokenize() {
		existing, closer, err := batch.Get(key)
		if err != nil && err != pebble.ErrNotFound {
			return err
		}

		var vCopy []byte
		if existing != nil {
			vCopy = slices.Clone(existing)
		}
		if closer != nil {
			closer.Close()
		}

		rb, err := decode(vCopy)
		if err != nil {
			return err
		}

		if !rb.Contains(numericID) {
			rb.Add(numericID)
		}

		blob, err := encode(rb)
		if err != nil {
			return err
		}

		if err = batch.Set(key, blob, pebble.Sync); err != nil {
			return err
		}

		if closer != nil {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateUser(c echo.Context, user *models.User) error {
	numID, err := GetOrAllocate(Client, user.GetID())
	if err != nil {
		return err
	}
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	keyUser := models.MakeKey("user_id", user.GetID())
	j, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if err = batch.Set(keyUser, j, pebble.Sync); err != nil {
		return err
	}

	if err = writeTokens(c, batch, user, numID); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// DeleteUser removes the user object and its posting-list entries.
func DeleteUser(c echo.Context, user *models.User) error {
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	ulid := user.GetID()

	numID, ok, err := Lookup(Client, ulid)
	if err != nil || !ok {
		return err // not found = nothing to delete
	}

	keyUser := models.MakeKey("user_id", ulid)
	oldBytes, closer, err := batch.Get(keyUser)
	if err == pebble.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	defer closer.Close()

    var oldUser models.User
	if err = json.Unmarshal(oldBytes, &oldUser); err != nil {
		return err
	}

	for _, tk := range oldUser.Tokenize() {
		data, closer, err := batch.Get(tk)
		if err != nil {
			return err
		}

		rb, err := decode(data)
		if err != nil {
			return err
		}
		if rb.Remove(numID); rb.IsEmpty() {
			if err = batch.Delete(tk, nil); err != nil {
				return err
			}
			if closer != nil {
				if err := closer.Close(); err != nil {
					return err
				}
			}
			continue
		}
		blob, _ := encode(rb)
		if err = batch.Set(tk, blob, pebble.Sync); err != nil {
			return err
		}

		if closer != nil {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}

	if err = batch.Delete(keyUser, nil); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// UpdateUser overwrites the user JSON and refreshes all bitmap tokens.
// Simplest strategy: delete old postings then re-add new ones.
func UpdateUser(c echo.Context, user *models.User) error {
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	ulid := user.GetID()
	numID, ok, err := Lookup(Client, ulid)
	if err != nil {
		return err
	}
	if !ok {
		return models.ErrUserNotFound
	}

	keyUser := models.MakeKey("user_id", ulid)
	prevBytes, closer, err := batch.Get(keyUser)
	if err == pebble.ErrNotFound {
		return models.ErrUserNotFound
	}
	if err != nil {
		return err
	}
	defer closer.Close()

    var prev models.User
	if err = json.Unmarshal(prevBytes, &prev); err != nil {
		return err
	}

	for _, tk := range prev.Tokenize() {
		data, closeFn, err := batch.Get(tk)
		if err != nil && err != pebble.ErrNotFound {
			return err
		}

		rb, err := decode(data)
		if err != nil {
			return err
		}
		if rb.Remove(numID); rb.IsEmpty() {
			if err = batch.Delete(tk, pebble.Sync); err != nil {
				return err
			}
			if closeFn != nil {
				if err := closeFn.Close(); err != nil {
					return err
				}
			}
			continue
		}
		blob, _ := encode(rb)
		if err = batch.Set(tk, blob, nil); err != nil {
			return err
		}

		if closeFn != nil {
			if err := closeFn.Close(); err != nil {
				return err
			}
		}
	}

	if err = writeTokens(c, batch, user, numID); err != nil {
		return err
	}

	newJSON, _ := json.Marshal(user)
	if err = batch.Set(keyUser, newJSON, pebble.Sync); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

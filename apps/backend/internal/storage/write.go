package storage

import (
	"encoding/json"
	"slices"

	"github.com/RoaringBitmap/roaring/roaring64"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/models"
	"github.com/Taiterbase/vtrips/apps/backend/pkg/utils"
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

func writeTokens(c echo.Context, batch *pebble.Batch, trip models.Trip, numericID uint64) (err error) {
	defer func() {
		if err != nil {
			c.Logger().Errorj(log.JSON{"err": err})
		}
	}()

	for _, key := range trip.Tokenize() {
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

func CreateTrip(c echo.Context, trip models.Trip) error {
	numID, err := GetOrAllocate(Client, trip.GetID())
	if err != nil {
		return err
	}
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	keyTrip := utils.MakeKey("trip_id", trip.GetID())
	j, err := json.Marshal(trip)
	if err != nil {
		return err
	}
	if err = batch.Set(keyTrip, j, pebble.Sync); err != nil {
		return err
	}

	if err = writeTokens(c, batch, trip, numID); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// DeleteTrip removes the trip object and its posting-list entries.
func DeleteTrip(c echo.Context, trip models.Trip) error {
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	ulid := trip.GetID()

	numID, ok, err := Lookup(Client, ulid)
	if err != nil || !ok {
		return err // not found = nothing to delete
	}

	keyTrip := utils.MakeKey("trip_id", ulid)
	oldBytes, closer, err := batch.Get(keyTrip)
	if err == pebble.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	defer closer.Close()

	var oldTrip models.TripBase
	if err = json.Unmarshal(oldBytes, &oldTrip); err != nil {
		return err
	}

	for _, tk := range oldTrip.Tokenize() {
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

	if err = batch.Delete(keyTrip, nil); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

// UpdateTrip overwrites the trip JSON and refreshes all bitmap tokens.
// Simplest strategy: delete old postings then re-add new ones.
func UpdateTrip(c echo.Context, trip models.Trip) error {
	batch := Client.NewIndexedBatch()
	defer batch.Close()

	ulid := trip.GetID()
	numID, ok, err := Lookup(Client, ulid)
	if err != nil {
		return err
	}
	if !ok {
		return models.ErrTripNotFound
	}

	keyTrip := utils.MakeKey("trip_id", ulid)
	prevBytes, closer, err := batch.Get(keyTrip)
	if err == pebble.ErrNotFound {
		return models.ErrTripNotFound
	}
	if err != nil {
		return err
	}
	defer closer.Close()

	var prev models.TripBase
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

	if err = writeTokens(c, batch, trip, numID); err != nil {
		return err
	}

	newJSON, _ := json.Marshal(trip)
	if err = batch.Set(keyTrip, newJSON, pebble.Sync); err != nil {
		return err
	}
	return batch.Commit(pebble.Sync)
}

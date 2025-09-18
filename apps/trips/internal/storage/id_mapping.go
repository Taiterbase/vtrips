package storage

import (
	"encoding/binary"
	"errors"

	"github.com/cockroachdb/pebble"
)

const (
	kCounter = "idmap/ctr" // 8-byte big-endian uint64
	kForward = "idmap/u/"  // ULID -> uint64
	kReverse = "idmap/i/"  // uint64 -> ULID
)

func getUint64(b []byte) uint64 {
	if len(b) == 8 {
		return binary.BigEndian.Uint64(b)
	}
	return 0
}

func putUint64(n uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], n)
	return buf[:]
}

func Lookup(db *pebble.DB, ulid string) (uint64, bool, error) {
	v, closer, err := db.Get([]byte(kForward + ulid))
	if err == pebble.ErrNotFound {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	defer closer.Close()
	return getUint64(v), true, nil
}

func GetOrAllocate(db *pebble.DB, ulid string) (uint64, error) {
	if v, closer, err := db.Get([]byte(kForward + ulid)); err == nil {
		defer closer.Close()
		return getUint64(v), nil
	} else if !errors.Is(err, pebble.ErrNotFound) {
		return 0, err
	}

	b := db.NewIndexedBatch()
	defer b.Close()
	curVal, _, err := b.Get([]byte(kCounter))
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return 0, err
	}
	next := getUint64(curVal) + 1

	if err = b.Set([]byte(kCounter), putUint64(next), nil); err != nil {
		return 0, err
	}
	if err = b.Set([]byte(kForward+ulid), putUint64(next), nil); err != nil {
		return 0, err
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], next)
	if err := b.Set(append([]byte(kReverse), buf[:]...), []byte(ulid), nil); err != nil {
		return 0, err
	}

	if err = b.Commit(pebble.Sync); err != nil {
		return 0, err
	}
	return next, nil
}

func Reverse(db *pebble.DB, id uint64) (ulid string, ok bool, err error) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], id)
	key := append([]byte(kReverse), buf[:]...)

	v, closer, err := db.Get(key)
	if err == pebble.ErrNotFound {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	defer closer.Close()
	return string(v), true, nil
}

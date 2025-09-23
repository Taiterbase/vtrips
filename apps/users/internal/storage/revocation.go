package storage

import (
    "encoding/binary"
    "fmt"
)

func revokedKey(userID string) []byte {
    return []byte(fmt.Sprintf("auth_revoked_before:%s", userID))
}

// SetRevokedBefore sets the UNIX timestamp before which tokens are revoked for a user.
func SetRevokedBefore(userID string, ts int64) error {
    var buf [8]byte
    binary.BigEndian.PutUint64(buf[:], uint64(ts))
    return Client.Set(revokedKey(userID), buf[:], nil)
}

// GetRevokedBefore returns the UNIX timestamp before which tokens are revoked for a user.
// If not set, returns 0 and nil error.
func GetRevokedBefore(userID string) (int64, error) {
    v, closer, err := Client.Get(revokedKey(userID))
    if err != nil {
        if err.Error() == "pebble: not found" {
            return 0, nil
        }
        return 0, err
    }
    defer closer.Close()
    if len(v) < 8 {
        return 0, nil
    }
    ts := int64(binary.BigEndian.Uint64(v[:8]))
    return ts, nil
}

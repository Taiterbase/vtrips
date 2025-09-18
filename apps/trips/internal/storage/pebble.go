package storage

import (
	"context"
	"log"

	"github.com/cockroachdb/pebble"
)

var Client *pebble.DB

func Initialize(ctx context.Context) {
	db, err := pebble.Open("/tmp/test.db", &pebble.Options{
		ErrorIfExists: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Pebble DB initialized")
	Client = db
}

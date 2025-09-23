package main

import (
	"context"

	"github.com/Taiterbase/vtrips/apps/trips/internal/api"
	"github.com/Taiterbase/vtrips/apps/trips/internal/storage"
)

func main() {
	ctx := context.Background()
	storage.Initialize(ctx)
	api.StartAPI()
}

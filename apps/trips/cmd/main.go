package main

import (
	"context"

	"github.com/Taiterbase/vtrips/apps/backend/internal/api"
	"github.com/Taiterbase/vtrips/apps/backend/internal/storage"
)

func main() {
	ctx := context.Background()
	storage.Initialize(ctx)
	api.StartAPI()
}

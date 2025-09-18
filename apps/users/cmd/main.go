package main

import (
	"context"

	"github.com/Taiterbase/vtrips/apps/users/internal/api"
	"github.com/Taiterbase/vtrips/apps/users/internal/storage"
)

func main() {
	ctx := context.Background()
	storage.Initialize(ctx)
	api.StartAPI()
}

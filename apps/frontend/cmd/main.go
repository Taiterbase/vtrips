package main

import (
	"context"
	"log"

	"github.com/Taiterbase/vtrips/apps/frontend/internal/api"
)

func main() {
	ctx := context.Background()
	srv := api.New()
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}

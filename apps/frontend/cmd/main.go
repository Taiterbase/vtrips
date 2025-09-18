package main

import (
	"context"
	"log"

	"github.com/Taiterbase/vtrips/apps/frontend/internal/server"
)

func main() {
	ctx := context.Background()
	srv := server.New()
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}

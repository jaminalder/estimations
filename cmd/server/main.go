package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/jaminalder/estimations/internal/adapters/http"
	"github.com/jaminalder/estimations/internal/adapters/idgen"
	"github.com/jaminalder/estimations/internal/adapters/memory"
	"github.com/jaminalder/estimations/internal/adapters/sse"
	"github.com/jaminalder/estimations/internal/app"
)

func main() {
	// Wire dependencies
	repo := memory.NewRoomRepo()
	ids := idgen.NewRandom(10, 8)
	hub := sse.NewHub(16)
	svc := &app.Service{Rooms: repo, Ids: ids, Bus: hub}

	// Renderer and server
	rend, err := httpadapter.NewRenderer()
	if err != nil {
		log.Fatalf("templates: %v", err)
	}
	handler := httpadapter.NewServer(svc, rend, httpadapter.WithLogger(log.Default()))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		log.Printf("HTTP listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

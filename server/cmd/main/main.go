package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ali-Karaki/e8markets/server/internal/config"
	"github.com/Ali-Karaki/e8markets/server/internal/handlers"
	"github.com/Ali-Karaki/e8markets/server/internal/middleware"
	"github.com/Ali-Karaki/e8markets/server/internal/server"
	"github.com/Ali-Karaki/e8markets/server/internal/store"
	"github.com/Ali-Karaki/e8markets/server/internal/tradelocker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	pg, err := store.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pg.Close()

	redisClient, err := store.NewRedis(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	sessions := store.NewSessionStore(redisClient)
	apiLogs := store.NewAPILogStore(pg)

	tl := tradelocker.NewClient(cfg.TradeLockerBaseURL, apiLogs)
	authMW := middleware.NewAuth(cfg, sessions, tl)

	authHandler := handlers.NewAuthHandler(cfg, sessions, tl, apiLogs, authMW)
	accountsHandler := handlers.NewAccountsHandler(tl)
	instrumentsHandler := handlers.NewInstrumentsHandler(tl)
	positionStore := store.NewPositionStore(pg)
	positionsHandler := handlers.NewPositionsHandler(tl, positionStore)

	handler := server.NewHandler(server.Deps{
		Cfg:         cfg,
		Auth:        authHandler,
		Accounts:    accountsHandler,
		Instruments: instrumentsHandler,
		Positions:   positionsHandler,
		AuthMW:      authMW,
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}

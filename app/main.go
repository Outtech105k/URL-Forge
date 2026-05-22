package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	redisclient "github.com/Outtech105k/ShortUrlServer/app/redis-client"
	"github.com/Outtech105k/ShortUrlServer/app/routes"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/caarlos0/env/v11"
)

func run() error {
	// Get Environment Variables
	var cfg utils.Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("getEnv: %w", err)
	}

	// Connct Redis
	redisAdapter, err := redisclient.NewRedisAdapter(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		return fmt.Errorf("connectRedis: %w", err)
	}
	defer redisAdapter.Close()

	// Migration
	if err := redisAdapter.MigrateAddPublicCtrl(false); err != nil {
		log.Printf("Migration warning: %v\n", err)
	}

	// Setup AppContext
	appCtx := &utils.AppContext{
		Config: cfg,
		Redis:  redisAdapter,
	}

	// Setup Gin Router
	router := routes.SetupRouter(appCtx)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	serverErrChan := make(chan error, 1)

	go func() {
		log.Printf("Starting server on :%s\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- fmt.Errorf("server listen error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit: // 終了信号検知
		log.Printf("Received signal: %s. Initiating shutdown...\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		log.Println("Server gracefully stopped")
		return nil

	case err := <-serverErrChan: // ListenAndServe中のエラー検知
		return err
	}
}

func main() {
	if err := run(); err != nil {
		log.Println("Application exited with error:", err)
		os.Exit(1)
	}
}

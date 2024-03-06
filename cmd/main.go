package main

import (
	"errors"
	"fmt"
	_ "github.com/zhayt/cert-tz/docs"
	"github.com/zhayt/cert-tz/internal/config"
	"github.com/zhayt/cert-tz/internal/service"
	"github.com/zhayt/cert-tz/internal/storage"
	"github.com/zhayt/cert-tz/internal/storage/postgre"
	cache "github.com/zhayt/cert-tz/internal/storage/redis"
	server "github.com/zhayt/cert-tz/internal/transport/http"
	"github.com/zhayt/cert-tz/internal/transport/http/handler"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// @title 	ЦАРКА REST API
// @version	1.0
// @description ЦАРКА Тествое задание.
// @termsOfService	http://swagger.io/terms/
// @host localhost:8000
// @BasePath	/rest
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var once sync.Once
	once.Do(config.PrepareENV)

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	l, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	defer func(*zap.Logger) {
		err := l.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(l)

	db, err := postgre.Dial(constructDSN(cfg))
	if err != nil {
		return err
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		return err
	}

	repo := storage.NewStorage(db, redisClient)

	usecases := service.NewService(repo, l)

	controller := handler.NewHandler(usecases, l)

	httpServer := server.NewServer(cfg, controller)

	l.Info("Start app", zap.String("port", cfg.AppPort))
	httpServer.StartServer()

	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err = <-httpServer.Notify:
		l.Info("server closing", zap.Error(err))
	}

	if err = httpServer.Shutdown(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error while shutting down server: %s", err)
	}

	return nil
}

func constructDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.TZ)
}

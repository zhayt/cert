package main

import (
	"fmt"
	"github.com/zhayt/cert-tz/config"
	"github.com/zhayt/cert-tz/service"
	"github.com/zhayt/cert-tz/storage"
	"github.com/zhayt/cert-tz/storage/postgre"
	server "github.com/zhayt/cert-tz/transport/http"
	"github.com/zhayt/cert-tz/transport/http/handler"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// load env
	var once sync.Once
	once.Do(config.PrepareENV)

	// get config
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	// init logger

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

	// init repo layer
	db, err := postgre.Dial(constructDSN(cfg))
	if err != nil {
		return err
	}
	defer db.Close()

	repo := storage.NewStorage(db)

	// init service layer
	usecases := service.NewService(repo, l)

	// init handler layer
	controller := handler.NewHandler(usecases, l)

	// init http server instance
	httpServer := server.NewServer(cfg, controller)

	l.Info("Start app", zap.String("port", cfg.AppPort))
	httpServer.StartServer()

	// grace full shutdown
	osSignCh := make(chan os.Signal, 1)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-osSignCh:
		l.Info("signal accepted: ", zap.String("signal", s.String()))
	case err = <-httpServer.Notify:
		l.Info("server closing", zap.Error(err))
	}

	if err = httpServer.Shutdown(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error while shutting down server: %s", err)
	}

	return nil
}

func constructDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.TZ)
}

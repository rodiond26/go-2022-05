package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/app"
	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	log.Printf("start\n")

	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	if err := config.ReadConfig(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	Password := os.Getenv("DB_PASSWORD")
	log.Printf("Password = [%v]\n", Password)

	logg := logger.New(config.Logger.Level)
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

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
	appConfig "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/config"
	appLogger "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
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

	mainConfig := appConfig.NewConfig()
	if err := mainConfig.ReadConfig(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	Password := os.Getenv("DB_PASSWORD")
	log.Printf("Password = [%v]\n", Password)

	logz, err := appLogger.NewLogger(&mainConfig)
	if err != nil {
		log.Fatal(err)
	}
	logz.Info("logz is Info")
	logz.Debug("logz is Debug")
	logz.Warn("logz is Warn")
	logz.Error("logz is Error")

	storage := memorystorage.New()
	calendar := app.New(logz, storage)
	server := internalhttp.NewServer(logz, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logz.Error("failed to stop http server: " + err.Error())
		}
	}()

	logz.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logz.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

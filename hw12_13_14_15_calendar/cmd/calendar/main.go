package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/app"
	appConfig "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/config"
	appLogger "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/server/http"
	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/initializing"
)

var configFile string

func init() {
	log.Printf("[0] initializing arguments ...\n")
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	log.Printf("[1] starting application ...\n")

	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	log.Printf("[2] initializing configuration ...\n")
	mainConfig := appConfig.NewConfig()
	if err := mainConfig.ReadConfig(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Printf("[3] initializing environment ...\n")
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal(err)
	// }

	log.Printf("[4] initializing logger ...\n")
	logz, err := appLogger.NewLogger(mainConfig.Environment.Type, mainConfig.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	log.Printf("[5] initializing storage ...\n")
	storage, err := storage.NewStorage(ctx, "", "dsn string")

	log.Printf("[6] initializing calendar ...\n")
	calendar := app.New(logz, storage)

	log.Printf("[7] initializing server ...\n")
	server := internalhttp.NewServer(logz, *calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Printf("[8] stopping server ...\n")
		if err := server.Stop(ctx); err != nil {
			logz.Error("failed to stop http server: " + err.Error())
		}
	}()

	logz.Info("calendar is running...")

	if err := server.Start(ctx, "localhost:8080"); err != nil {
		logz.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

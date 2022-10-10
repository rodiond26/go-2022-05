package main

import (
	"context"
	"flag"
	"fmt"
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
	storage "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/storage/init_storage"
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
	log.Printf("Password = [%v]\n", Password) // TODO delete

	logz, err := appLogger.NewLogger(mainConfig.Environment.Type, mainConfig.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}
	logz.Info("logz is Info")
	logz.Debug("logz is Debug")
	logz.Warn("logz is Warn")
	logz.Error("logz is Error")

	ctx := context.Background()

	storage, err := storage.NewStorage(ctx, "", "dsn string")
	fmt.Printf("storage = %+v]\n", storage)
	calendar := app.New(logz, storage)
	fmt.Printf("calendar = %+v]\n", calendar)

	server := internalhttp.NewServer(logz)
	fmt.Printf("server = %+v]\n", server)

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

	if err := server.Start(ctx, "localhost:8080"); err != nil {
		logz.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

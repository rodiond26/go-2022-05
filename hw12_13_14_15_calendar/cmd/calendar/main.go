package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/app"
	appConfig "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/config"
	appLogger "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/rodiond26/go-2022-05/hw12_13_14_15_calendar/internal/server/grpc"
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
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[6] initializing calendar ...\n")
	calendar := app.New(logz, storage)

	log.Printf("[7] initializing http server ...\n")
	httpServer := internalhttp.NewServer(logz, calendar)

	log.Printf("[8] initializing grpc server ...\n")
	grpcServer := internalgrpc.NewServer(logz, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Printf("[9] stopping server ...\n")
		if err := httpServer.Stop(ctx); err != nil {
			logz.Error("failed to stop http server: " + err.Error())
		}

		if err = grpcServer.Stop(ctx); err != nil {
			logz.Error("failed to stop grpc server: " + err.Error())
		}

		if err = calendar.Close(ctx); err != nil {
			logz.Error("failed close storage: " + err.Error())
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(mainConfig.Server.Host, mainConfig.Server.HttpPort)
		if err := httpServer.Start(ctx, addrServer); err != nil {
			logz.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	go func() {
		defer wg.Done()
		addrServer := net.JoinHostPort(mainConfig.Server.Host, mainConfig.Server.GrpcPort)
		if err := grpcServer.Start(ctx, addrServer); err != nil {
			logz.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}
	}()

	<-ctx.Done()
	wg.Wait()
}

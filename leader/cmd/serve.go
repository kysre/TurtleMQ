package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kysre/TurtleMQ/leader/internal/pkg/grpcserver"

	"github.com/kysre/TurtleMQ/leader/internal/app/core"
	"github.com/kysre/TurtleMQ/leader/pkg/queue"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start Server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	serverCtx, serverCancel := makeServerCtx()
	defer serverCancel()

	config, err := provideConfig(cmd)
	if err != nil {
		panic(err)
	}

	logger := getLoggerOrPanic(config)
	server := getServerOrPanic(config, logger)

	var serverWaitGroup sync.WaitGroup
	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()

		if err := server.Serve(); err != nil {
			panicWithError(err, "failed to serve")
		}
	}()

	if err := declareReadiness(); err != nil {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	server.Stop()

	serverWaitGroup.Wait()
}

func getLoggerOrPanic(conf *Config) *logrus.Logger {
	logger, err := provideLogger(conf)
	if err != nil {
		panic(err)
	}
	return logger
}

func getServerOrPanic(conf *Config, log *logrus.Logger) *grpcserver.Server {
	server, err := provideServer(core.New(), conf, log)
	if err != nil {
		panic(err)
	}
	return server
}

func provideServer(server queue.QueueServer, config *Config, logger *logrus.Logger) (*grpcserver.Server, error) {
	return grpcserver.New(server, logger, config.GRPC.ListenPort)
}

func makeServerCtx() (context.Context, context.CancelFunc) {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-gracefulStop
		cancel()
	}()

	return ctx, cancel
}

func declareReadiness() error {
	// nolint: gosec
	file, err := os.Create("/tmp/readiness")
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer file.Close()

	_, err = file.WriteString("ready")
	return err
}

func panicWithError(err error, format string, args ...interface{}) {
	logrus.WithError(err).Panicf(format, args...)
}
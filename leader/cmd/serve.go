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

	"github.com/kysre/TurtleMQ/leader/cmd/tasks"
	"github.com/kysre/TurtleMQ/leader/internal/app/core"
	"github.com/kysre/TurtleMQ/leader/internal/app/loadbalancer"
	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/internal/pkg/grpcserver"
	"github.com/kysre/TurtleMQ/leader/internal/pkg/metrics/prometheus"
	"github.com/kysre/TurtleMQ/leader/pkg/leader"
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

	// Start metrics server
	prometheusMetricServer := prometheus.StartMetricServerOrPanic(config.Metric.ListenPort)
	defer prometheus.ShutDownMetricsServer(prometheusMetricServer)

	// Get shared resources
	logger := getLoggerOrPanic(config)
	directory := getDataNodeDirectoryOrPanic()
	balancer := getLoadBalancerOrPanic(logger, directory)
	runDataNodesTasks(config, logger, directory, balancer)

	// Get grpc server cores
	queueCore := getQueueCore(logger, directory, balancer)
	leaderCore := getLeaderCore(logger, directory, balancer)
	// Register grpc server
	queueServer := getServerOrPanic(config, logger, queueCore, leaderCore)
	// Serve
	var serverWaitGroup sync.WaitGroup
	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()

		logger.Info("Start serving...")

		if err := queueServer.Serve(); err != nil {
			panicWithError(err, "failed to serve")
		}

		logger.Info("End serving...")
	}()

	if err := declareReadiness(); err != nil {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	queueServer.Stop()

	serverWaitGroup.Wait()
}

func getLoggerOrPanic(conf *Config) *logrus.Logger {
	logger, err := provideLogger(conf)
	if err != nil {
		panic(err)
	}
	return logger
}

func getDataNodeDirectoryOrPanic() *models.DataNodeDirectory {
	directory := models.NewDataNodeDirectory()
	if directory == nil {
		panic("DataNodeDirectory is nil")
	}
	return directory
}

func getLoadBalancerOrPanic(log *logrus.Logger, directory *models.DataNodeDirectory) loadbalancer.Balancer {
	return loadbalancer.NewBalancer(log, directory)
}

func runDataNodesTasks(
	conf *Config, log *logrus.Logger, directory *models.DataNodeDirectory, balancer loadbalancer.Balancer,
) {
	syncer := tasks.NewDataNodeSyncer(
		log, directory, balancer, conf.Leader.DataNodePartitionCount, conf.Leader.DataNodeSyncTimeout)
	healthChecker := tasks.NewDataNodeHealthChecker(directory, conf.Leader.DataNodeStateCheckPeriod, syncer)
	go healthChecker.RunHealthChecks()
	go tasks.RunRemainingCheck(directory, conf.Leader.DataNodeRemainingCheckPeriod)
}

func getQueueCore(
	log *logrus.Logger, directory *models.DataNodeDirectory, balancer loadbalancer.Balancer,
) queue.QueueServer {
	return core.NewQueueCore(log, directory, balancer)
}

func getLeaderCore(
	log *logrus.Logger, directory *models.DataNodeDirectory, balancer loadbalancer.Balancer,
) leader.LeaderServer {
	return core.NewLeaderCore(log, directory, balancer)
}

func getServerOrPanic(
	conf *Config, log *logrus.Logger, queueCore queue.QueueServer, leaderCore leader.LeaderServer,
) *grpcserver.Server {
	server, err := provideServer(conf, log, queueCore, leaderCore)
	if err != nil {
		panic(err)
	}
	return server
}

func provideServer(
	config *Config, logger *logrus.Logger, queueServicer queue.QueueServer, leaderServicer leader.LeaderServer,
) (*grpcserver.Server, error) {
	baseServer, err := grpcserver.New(logger, config.Queue.ListenPort)
	if err != nil {
		return nil, err
	}
	baseServer.RegisterQueueServer(queueServicer)
	logger.Info("Registered Queue Server")
	baseServer.RegisterLeaderServer(leaderServicer)
	logger.Info("Registered Leader Server")
	return baseServer, nil
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

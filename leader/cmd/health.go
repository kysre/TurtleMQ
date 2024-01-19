package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "health-check the server",
	Run:   healthCheck,
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func healthCheck(cmd *cobra.Command, _ []string) {
	config, err := LoadConfig(cmd)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%d", config.Queue.ListenPort)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("failed to dial grpc server")
		os.Exit(1)
	}

	client := grpc_health_v1.NewHealthClient(conn)
	if _, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{}); err != nil {
		logrus.WithError(err).Error("failed to health check server")
		os.Exit(1)
	}
}

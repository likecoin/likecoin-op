package cmd

import (
	"context"
	"os"

	"github.com/likecoin/like-migration-backend/cmd/worker/cmd/enqueue"
	"github.com/likecoin/like-migration-backend/cmd/worker/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Migration Backend worker CLI",
	Long:  `CLI to perform worker operations`,
}

func Execute(envCfg *config.EnvConfig) {
	ctx := context.WithValue(context.Background(), config.ContextKey, envCfg)
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(WorkerCmd)
	rootCmd.AddCommand(SchedulerCmd)
	rootCmd.AddCommand(enqueue.EnqueueCmd)
}

package cmd

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"

	"github.com/likecoin/like-migration-backend/cmd/worker/context"
	"github.com/likecoin/like-migration-backend/cmd/worker/task"
)

var SchedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Start scheduelr",
	Run: func(cmd *cobra.Command, args []string) {
		envCfg := context.ConfigFromContext(cmd.Context())

		opt, err := redis.ParseURL(envCfg.RedisDsn)
		if err != nil {
			log.Fatalf("could not parse redis url: %v", err)
		}
		scheduler := asynq.NewScheduler(
			asynq.RedisClientOpt{
				Network:      opt.Network,
				Addr:         opt.Addr,
				Password:     opt.Password,
				DB:           opt.DB,
				DialTimeout:  opt.DialTimeout,
				ReadTimeout:  opt.ReadTimeout,
				WriteTimeout: opt.WriteTimeout,
				PoolSize:     opt.PoolSize,
				TLSConfig:    opt.TLSConfig,
			}, nil)
		task, err := task.NewHelloWorldTask("periodic")
		if err != nil {
			log.Fatalf("could not create task: %v", err)
		}

		// ... Register tasks
		_, err = scheduler.Register("* * * * *", task)

		if err != nil {
			log.Fatalf("could not register task: %v", err)
		}

		if err := scheduler.Run(); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	},
}

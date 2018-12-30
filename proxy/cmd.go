package proxy

import (
	"github.com/spf13/cobra"
	"sync"
)

func GetCommand() *cobra.Command {
	var redisAddress string
	var maxItems int64

	cmd := &cobra.Command{
		Use: "proxy",
		Short: "Starts a proxy server on the public network to listen for web hooks.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			redis, err := connectToRedis(redisAddress)
			if err != nil {
				return err
			}

			server := &server{
				redis: redis,
				pushLock: &sync.Mutex{},
				popLock: &sync.Mutex{},
				maxItems: maxItems,
			}
			return server.startServer(8080)
		},
	}

	cmd.PersistentFlags().StringVarP(&redisAddress, "redis-address", "r", "localhost:6379", "Address to the Redis database.")
	cmd.PersistentFlags().Int64VarP(&maxItems, "max-items", "x", 500, "Maximum number of events to store in Redis before starting to drop the oldest.")

	return cmd
}
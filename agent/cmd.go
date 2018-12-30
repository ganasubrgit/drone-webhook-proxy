package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	var proxyUrl, droneHookUrl string
	var interval int64

	cmd := &cobra.Command{
		Use: "agent",
		Short: "Starts a agent along side private drone deployment to poll proxy for saved web hooks.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infof("Polling remote proxy at %s and deliver to local drone at %s.", proxyUrl, droneHookUrl)

			ticker := &ticker{
				interval: interval,
				proxyUrl: proxyUrl,
				droneHookUrl: droneHookUrl,
			}

			never := make(chan struct{})
			ticker.start()
			<-never

			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&proxyUrl, "proxy-url", "p", "", "URL to the pop API of the remote proxy instance.")
	cmd.PersistentFlags().StringVarP(&droneHookUrl, "drone-hook-url", "d", "", "URL to the hook API of the local drone instance.")
	cmd.PersistentFlags().Int64VarP(&interval, "interval", "i", 5, "Frequency in seconds to poll remote proxy.")

	return cmd
}
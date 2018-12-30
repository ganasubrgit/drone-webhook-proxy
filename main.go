package main

import (
	"flag"
	"github.com/imulab/drone-webhook-proxy/agent"
	"github.com/imulab/drone-webhook-proxy/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var rootCommand = &cobra.Command{
	Use: "hook",
	Short: "Web hook proxy for drone deployments in private network.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.ParseFlags(args); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	var debug bool

	rootCommand.AddCommand(proxy.GetCommand())
	rootCommand.AddCommand(agent.GetCommand())
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Sets the log output level to debug.")
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if err := rootCommand.Execute(); err != nil {
		logrus.Errorf("Error running command.", err)
		os.Exit(1)
	}
}

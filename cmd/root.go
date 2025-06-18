package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/StevenWeathers/hatch-messaging-service/internal/config"

	"github.com/spf13/cobra"
)

var cfgFile string
var c *config.Config

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hms",
	Short: "Example hatch messaging service",
	Long: `To get started run the serve subcommand which will start a server
on localhost:8080:

    hms serve
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		slog.Error("error executing root command", slog.Any("error", err))
		os.Exit(-1)
	}
}

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hms.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	c = config.New(context.Background(), cfgFile)
}

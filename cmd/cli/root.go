package cli

import (
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/logging"
	"os"

	"github.com/spf13/cobra"
)

var settings = config.GetSettings()
var logger = logging.GetRootLogger(settings.Logging)

var rootCommand = &cobra.Command{
	Use:   "app",
	Short: "app",
	Long:  "app",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("No-op")
	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		logger.Error("Error executing CLI", "error", err)
		os.Exit(1)
	}
}

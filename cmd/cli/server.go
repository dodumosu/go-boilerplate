package cli

import (
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/logging"
	"go-boilerplate/internal/server"
	"os"

	"github.com/spf13/cobra"
)

func runServer() {
	settings := config.GetSettings()
	rootLogger := logging.GetRootLogger(settings.Logging)

	server, err := server.NewServer(rootLogger, settings)
	if err != nil {
		rootLogger.Error("Error creating server", "error", err)
		os.Exit(1)
	}
	server.Start()
}

var serverCommand = &cobra.Command{
	Use:     "run",
	Aliases: []string{"server"},
	Short:   "Runs the server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCommand.AddCommand(serverCommand)
}

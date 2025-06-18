package cmd

import (
	"context"

	"github.com/StevenWeathers/hatch-messaging-service/internal/db"
	"github.com/StevenWeathers/hatch-messaging-service/internal/http"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches the service on https://localhost:8080",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serve() {
	_, err := db.New(db.Config{
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		User:     c.Database.User,
		Password: c.Database.Password,
		DBName:   c.Database.Name,
	})
	if err != nil {
		panic(err)
	}

	h := http.New(http.Config{
		ListenAddress: c.ListenAddress,
	})

	err = h.ListenAndServe(context.Background())
	if err != nil {
		panic(err)
	}
}

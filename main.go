package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/fabzo/gcloud-directory-service/cmd/server"
)

var RootCmd = &cobra.Command{
	Use:   "gcloud-directory-service",
	Short: "Server that provides cached access to the google groups directory",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCmd.AddCommand(server.Command)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
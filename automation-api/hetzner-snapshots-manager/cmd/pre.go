package cmd

import (
	"automation-api/hetzner-snapshots-manager/manager"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var preCmd = &cobra.Command{
	Use:   "pre",
	Short: "The 'pre' only make preview",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		m, err := manager.New(ctx)
		if err != nil {
			log.Fatalf("create manager: %s", err)
		}

		if err := m.Run(cmd, true); err != nil {
			m.Logger.Fatal(fmt.Sprintf("run: %s", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(preCmd)
}

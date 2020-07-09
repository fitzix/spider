package cmd

import (
	"github.com/fitzix/spider/whatsapp"
	"github.com/spf13/cobra"
)

var chromePath string

func init() {
	rootCmd.AddCommand(whatsAppCmd)
	whatsAppCmd.Flags().StringVarP(&chromePath, "chrome", "c", "", "chrome app path")
}

var whatsAppCmd = &cobra.Command{
	Use: "whatsapp",
	Run: func(cmd *cobra.Command, args []string) {
		whatsapp.Run()
	},
}

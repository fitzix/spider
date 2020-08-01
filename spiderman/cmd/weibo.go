package cmd

import (
	"github.com/fitzix/spider/services"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(weiboCmd)
}

var weiboCmd = &cobra.Command{
	Use: "weibo",
	Run: func(cmd *cobra.Command, args []string) {
		services.Weibo()
	},
}

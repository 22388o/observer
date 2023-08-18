/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"obwallet/obrpc/user"

	"github.com/spf13/cobra"
)

// vcodeCmd represents the vcode command
var vcodeCmd = &cobra.Command{
	Use:   "vcode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		in := &user.VerifyCodeRequest{Email: "wxf4150@163.com"}
		res, err := gClient.VerifyCode(context.TODO(), in)
		if err != nil {
			log.Println(err)
			return
		}
		printRespJSON(res)
	},
}

func init() {
	rootCmd.AddCommand(vcodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vcodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vcodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

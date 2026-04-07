package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "resource-template-coverter",
	Short: "Tool for migrating resource yaml templates",
	Long:  `Tool for migrating resource yaml templates`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

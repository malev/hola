package cmd

import (
	"fmt"

	"github.com/malev/hola/internals"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hola",
	Long:  `All software has versions. This is Hola's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hola HTTP Client " + internals.Version)
	},
}

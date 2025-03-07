package cmd

import (
	"fmt"
	"os"

	"github.com/Arthx-x/arthxrecon/util"
	"github.com/spf13/cobra"
)

// rootCmd is the main command for the application.
var rootCmd = &cobra.Command{
	Use:   util.AppName,
	Short: util.AppDescription,
	Long:  util.AppDescription,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.Banner() // Chama seu banner antes de qualquer comando ser executado
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.CmdUsage)
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Registra os subcomandos
	rootCmd.AddCommand(HostDiscoveryCmd)
	rootCmd.AddCommand(PortScanCmd)
	// Você pode adicionar outros subcomandos, como portscan, enumeration, etc.
}

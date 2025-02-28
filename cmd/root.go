package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd é o comando principal da aplicação.
var rootCmd = &cobra.Command{
	Use:   "arthxrecon",
	Short: "ArthxRecon is a modular recon tool for pentesting",
	Long:  "ArthxRecon is a modular network reconnaissance tool that integrates multiple scanning and enumeration modules.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ArthxRecon: please specify a subcommand (e.g., portscan, hostdiscovery, etc.)")
	},
}

// Execute executa o comando raiz.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

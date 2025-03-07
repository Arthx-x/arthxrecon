package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/Arthx-x/arthxrecon/internal/hostdiscovery"
	"github.com/Arthx-x/arthxrecon/util"
	"github.com/spf13/cobra"
)

var (
	hostTarget        string // Alvo único (IP ou CIDR) passado via flag
	hostOutputFile    string // Nome base para os arquivos de saída
	hostMode          string // Modo do scan: aggressive, normal ou passive
	hostCustomOptions string
)

// HostDiscoveryCmd é o comando para executar a descoberta de hosts.
var HostDiscoveryCmd = &cobra.Command{
	Use:   util.HostDiscoveryName,
	Short: util.HDAppDescription,
	Run: func(cmd *cobra.Command, args []string) {
		// Decide qual alvo utilizar, seja da flag ou do arquivo
		targets, fileMode := util.ParseTargetInput(hostTarget)
		if len(targets) == 0 {
			log.Fatal().Msgf("%s", util.ErrInvalidTarget)
		}
		options := []string{}
		if hostCustomOptions != "" {
			options = strings.Fields(hostCustomOptions)
			// Opcional: faça trim em cada opção
			for i, opt := range options {
				options[i] = strings.TrimSpace(opt)
			}
		}

		// Configura os parâmetros para a descoberta de hosts.
		params := hostdiscovery.DiscoveryParams{
			Targets:    targets,        // Lista de alvos.
			OutputFile: hostOutputFile, // Nome base para os arquivos de saída.
			Mode:       hostMode,       // Modo do scan.
			Options:    options,        // Outras opções extras, se necessário.
			FileMode:   fileMode,       // Indica se os targets vieram de um arquivo.
		}

		// Seleciona a estratégia de descoberta (aqui usamos o Nmap como padrão)
		strategy := hostdiscovery.NewNmapHostDiscovery()

		// Cria o orquestrador que gerencia o fluxo completo: configurar, executar e parsear
		orchestrator := hostdiscovery.NewHostDiscoveryOrchestrator(strategy, params)

		fmt.Printf("%s Host Discovery", util.MarkerCyan)
		fmt.Printf("\n%s %s Starting\n", util.MarkerCyan, util.GetFormattedTime())
		hosts, err := orchestrator.Run()

		if err != nil {
			log.Fatal().Msgf("%s %v", util.FatalErrHD, err)
		}
		fmt.Printf("%s Discovered: %s %s\n", util.MarkerGreen, util.Green(strconv.Itoa(len(hosts))), util.Green("Hosts"))
		fmt.Printf("\n%s %s Finished\n", util.MarkerCyan, util.GetFormattedTime())
	},
}

func init() {
	HostDiscoveryCmd.Flags().StringVarP(&hostTarget, "target", "t", "", "Target IP(s) or CIDR range, or path to file containing targets (if file, provide file path; for multiple, separate by commas)")
	HostDiscoveryCmd.Flags().StringVarP(&hostOutputFile, "outfile", "o", "targets", "Base name for output files")
	HostDiscoveryCmd.Flags().StringVarP(&hostMode, "mode", "m", "normal", "Scan mode: 1.stealth, 2.normal, or 3.aggressive")
	HostDiscoveryCmd.Flags().StringVarP(&hostCustomOptions, "custom", "c", "", "Custom options for the scan, separated by commas")
	// Adicione o comando ao rootCmd em root.go.
}

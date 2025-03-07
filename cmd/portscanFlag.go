package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/Arthx-x/arthxrecon/internal/portscan"
	"github.com/Arthx-x/arthxrecon/util"
	"github.com/spf13/cobra"
)

var (
	psTarget        string // Alvo único ou múltiplos (IP(s) ou CIDR) passado via flag ou arquivo
	psOutputFile    string // Nome base para os arquivos de saída
	psPortList      string // Lista ou range de portas (ex: "1-1024")
	psMode          string // Modo do scan: aggressive, normal ou passive
	psCategory      string // Categorias de portas a incluir (ex: top12, database, web, network, firewall, windows, vpn, all)
	psAllPorts      bool   // Se definido, varre todas as portas (-p-) em execução separada em background
	psSimpleScan    bool   // Se definido, realiza um portScan simples (ex.: -sS); caso contrário, usa -sV -sC
	psCustomOptions string // Opções customizadas extras para o scan, separadas por espaços
)

// PortScanCmd é o comando para executar a varredura de portas.
var PortScanCmd = &cobra.Command{
	Use:   "portscan",
	Short: "Performs a port scan on specified targets",
	Run: func(cmd *cobra.Command, args []string) {
		// Processa os alvos (pode ser via flag ou arquivo)
		targets, fileMode := util.ParseTargetInput(psTarget)
		if len(targets) == 0 {
			log.Fatal().Msg(util.ErrInvalidTarget)
		}

		// Processa as opções customizadas, convertendo a string para um slice de strings.
		options := []string{}
		if psCustomOptions != "" {
			options = strings.Fields(psCustomOptions)
			for i, opt := range options {
				options[i] = strings.TrimSpace(opt)
			}
		}

		// Configura os parâmetros para o port scan.
		params := portscan.PortScanParams{
			Targets:    targets,      // Lista de alvos.
			OutputFile: psOutputFile, // Nome base para os arquivos de saída.
			Mode:       psMode,       // Modo do scan.
			Options:    options,      // Opções customizadas extras.
			PortList:   psPortList,   // Lista ou range de portas.
			Category:   psCategory,   // Categoria de portas, se definida.
			AllPorts:   psAllPorts,   // Flag para varredura de todas as portas.
			SimpleScan: psSimpleScan, // Flag para usar um scan simples.
			FileMode:   fileMode,     // Indica se os alvos vieram de um arquivo.
		}

		// Seleciona a estratégia de port scan (aqui usamos Nmap como padrão).
		strategy := portscan.NewNmapPortScanner()

		// Exibe as configurações utilizadas.
		fmt.Printf("\n%s Port Scan", util.MarkerCyan)
		fmt.Printf("\n%s %s Starting\n", util.MarkerCyan, util.GetFormattedTime())

		// Cria o orquestrador para o port scan.
		orchestrator := portscan.NewPortScanOrchestrator(strategy, params)

		// Exibe as configurações uma única vez.
		portscan.ShowConfiguration(params)

		// Executa o scan.
		ports, err := orchestrator.Run()
		if err != nil {
			log.Fatal().Msgf("%s %v", util.FatalErrPS, err)
		}

		// Exibe o resultado (por exemplo, quantidade de portas descobertas).
		fmt.Printf("%s Ports discovered 01: %s\n", util.MarkerGreen, util.Green(strconv.Itoa(len(ports))))
		fmt.Printf("\n%s %s Finished\n", util.MarkerCyan, util.GetFormattedTime())
	},
}

func init() {
	PortScanCmd.Flags().StringVarP(&psTarget, "target", "t", "./hostDiscovery/targets.txt", "Target IP(s) or CIDR range, or path to file with targets (if file, provide path; for multiple, separate by commas)")
	PortScanCmd.Flags().StringVarP(&psOutputFile, "outfile", "o", "portscan", "Base name for output files")
	PortScanCmd.Flags().StringVarP(&psPortList, "ports", "p", "", "Port range or list to scan (e.g., \"1-1024\")")
	PortScanCmd.Flags().StringVarP(&psMode, "mode", "m", "normal", "Scan mode: aggressive, normal, or passive")
	PortScanCmd.Flags().StringVarP(&psCategory, "category", "c", "", "Port category to include (e.g., top12, database, web, network, firewall, windows, vpn, all)")
	PortScanCmd.Flags().BoolVarP(&psAllPorts, "allports", "a", false, "Scan all ports (-p-) in background")
	PortScanCmd.Flags().BoolVarP(&psSimpleScan, "simple", "s", false, "Use a simple port scan (e.g., -sS) instead of a detailed scan (-sV -sC)")
	PortScanCmd.Flags().StringVarP(&psCustomOptions, "custom", "x", "", "Custom options for the scan, separated by spaces")
}

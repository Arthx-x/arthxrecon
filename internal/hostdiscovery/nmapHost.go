package hostdiscovery

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Arthx-x/arthxrecon/util"
	"github.com/rs/zerolog/log"
	"github.com/tomsteele/go-nmap"
)

// NmapHostDiscovery implementa HostDiscoveryStrategy utilizando o Nmap.
type NmapHostDiscovery struct {
	Targets    []string // Lista de alvos (IPs ou CIDRs)
	OutputFile string   // Nome base para os arquivos de saída
	Mode       string   // Modo do scan (aggressive, normal, passive)
	Options    []string // Outras opções, se houver
	FileMode   bool     // Indica se os targets vieram de um arquivo (modo arquivo)
}

// NewNmapHostDiscovery é a factory que cria uma instância de NmapHostDiscovery.
func NewNmapHostDiscovery() *NmapHostDiscovery {
	return &NmapHostDiscovery{}
}

// Configure configura a estratégia com os parâmetros fornecidos.
func (nmapHD *NmapHostDiscovery) Configure(params DiscoveryParams) error {
	nmapHD.Targets = params.Targets
	nmapHD.OutputFile = filepath.Join(util.HostDiscoveryName, params.OutputFile)
	nmapHD.Mode = params.Mode
	nmapHD.Options = params.Options
	nmapHD.FileMode = params.FileMode

	// fmt.Printf("\n┌──────────────────────────────────────────────┐\n  %s Target \t: %s \n  %s Output \t: %s \n  %s Mode \t: %s \n  %s Options \t: %s \n└──────────────────────────────────────────────┘\n\n",
	// 	util.MarkerGreen,
	// 	strings.Join(nmapHD.Targets, ", "),
	// 	util.MarkerGreen,
	// 	nmapHD.OutputFile,
	// 	util.MarkerGreen,
	// 	nmapHD.Mode,
	// 	util.MarkerGreen,
	// 	strings.Join(nmapHD.Options, ", "),
	// )

	fmt.Printf("\n┌──────────────────────────────────────────────┐\n  %s \t: %s \n  %s \t: %s \n  %s \t: %s \n  %s \t: %s \n└──────────────────────────────────────────────┘\n\n",
		util.Green("⦿ Target"),
		strings.Join(nmapHD.Targets, ", "),
		util.Green("⦿ Output"),
		nmapHD.OutputFile,
		util.Green("⦿ Mode"),
		nmapHD.Mode,
		util.Green("⦿ Options"),
		strings.Join(nmapHD.Options, ", "),
	)

	//log.Debug().Msgf("Configure: Targets=%v, OutputFile=%s, Mode=%s, Options=%v, FileMode=%t",
	//	nmapHD.Targets, nmapHD.OutputFile, nmapHD.Mode, nmapHD.Options, nmapHD.FileMode)

	return nil
}

// buildCommand monta o comando nmap com base na configuração.
func (nmapHD *NmapHostDiscovery) buildCommand() (string, []string) {
	// Comando base: nmap -sn [opções] <alvo(s)> -oA <outputFile>

	outputDir := util.HostDiscoveryName
	if err := util.EnsureDir(outputDir); err != nil {
		log.Fatal().Msgf("Error creating directory %s: %v", outputDir, err)
	}

	args := []string{"-sn", util.HostDiscoveryFlagNmap}
	if len(nmapHD.Options) > 0 {
		args = append(args, nmapHD.Options...)
	}

	// Adiciona o flag de modo.
	modeFlag := ""
	lowerMode := strings.ToLower(strings.TrimSpace(nmapHD.Mode))
	switch lowerMode {
	case "aggressive", "3":
		modeFlag = "-T4"
	case "stealth", "1":
		modeFlag = "-T2"
		// Para stealth, pode-se adicionar opções extras se necessário.
	case "normal", "2":
		// Sem flag adicional.
	default:
		// Valor padrão.
	}
	if modeFlag != "" {
		args = append(args, modeFlag)
	}

	// Se for fileMode, utiliza o primeiro (único) target como caminho para o arquivo.
	if nmapHD.FileMode {
		args = append(args, "-iL", nmapHD.Targets[0])
	} else {
		// Caso contrário, adicione cada target individualmente.
		args = append(args, nmapHD.Targets...)
	}
	// define o nome para o arquivo de saida
	args = append(args, "-oA", nmapHD.OutputFile)
	commandStr := "nmap " + strings.Join(args, " ")
	return commandStr, args
}

// Execute executa o comando nmap e retorna a saída bruta.
func (nmapHD *NmapHostDiscovery) Execute() (string, error) {
	commandStr, args := nmapHD.buildCommand()
	//fmt.Printf("%s Running: ", util.MarkerGreen+util.Red(commandStr))
	fmt.Printf("%s Running: %s\n", util.MarkerGreen, util.Green(commandStr))
	//log.Info().Msgf("Executing host discovery: %s", commandStr)

	cmd := exec.Command("nmap", args...)
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("nmap execution failed: %w", err)
	}

	// Após a execução, lemos o arquivo XML gerado.
	xmlFilePath := nmapHD.OutputFile + ".xml"
	data, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read XML file: %w", err)
	}
	return string(data), nil
}

// Parse utiliza a biblioteca go-nmap para converter o XML em estrutura e extrair os hosts.
func (nmapHD *NmapHostDiscovery) Parse(rawOutput string) ([]string, error) {
	result, err := nmap.Parse([]byte(rawOutput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse nmap XML: %w", err)
	}
	var hosts []string

	for _, host := range result.Hosts {
		for _, address := range host.Addresses {
			if address.AddrType == "ipv4" {
				hosts = append(hosts, address.Addr)
				break
			}
		}
	}

	// Cria ou sobrescreve o arquivo targets.txt na pasta "hostDiscovery".
	targetsFile := filepath.Join(util.HostDiscoveryName, "targets.txt")
	if err := util.WriteTargetsToFile(targetsFile, hosts); err != nil {
		log.Error().Err(err).Msg("Failed to write targets.txt")
	} else {
		fmt.Printf("%s Creating: %s\n", util.MarkerGreen, util.Green(targetsFile))
		//log.Info().Msgf("Targets extracted successfully to: %s", targetsFile)
	}

	return hosts, nil
}

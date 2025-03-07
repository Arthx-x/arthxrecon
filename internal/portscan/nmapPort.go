package portscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arthx-x/arthxrecon/util"
	"github.com/rs/zerolog/log"
)

// PortScanParams centraliza os parâmetros para o port scan.
type PortScanParams struct {
	Targets    []string // Lista de alvos (IPs ou CIDRs)
	OutputFile string   // Nome base para os arquivos de saída
	Mode       string   // Modo do scan: aggressive, normal ou passive
	Options    []string // Outras opções extras para o scan
	PortList   string   // Lista ou range de portas (ex.: "1-1024")
	Category   string   // Categoria de portas (ex.: "top12", "database", etc.)
	AllPorts   bool     // Se verdadeiro, varre todas as portas (-p-)
	SimpleScan bool     // Se verdadeiro, usa -sS; caso contrário, usa -sV -sC
	FileMode   bool     // Se os alvos foram passados via arquivo
}

// NmapPortScanner implementa a interface PortScanStrategy usando Nmap.
type NmapPortScanner struct {
	Targets    []string
	OutputFile string
	Mode       string
	Options    []string
	PortList   string
	Category   string
	AllPorts   bool
	SimpleScan bool
	FileMode   bool
}

// NewNmapPortScanner é a factory que cria uma instância de NmapPortScanner.
func NewNmapPortScanner() *NmapPortScanner {
	return &NmapPortScanner{}
}

var portCategories = map[string]string{
	"top12":    "21,22,2222,23,53,80,135,139,443,445,3389,8080",
	"database": "3306,5432,1433,1521,27017,6379,9042,9160,50000,8086,5984,7474,7687,11211,3050,9092,1527,2638,8529,28015,2424,26257,9200",
	"web":      "80,443,8080,8443,8000,3000,5000,4200,8888,8081,8001,3001,9000,9090,1313,8008,8880",
	"network":  "10000,20000,902,903,8006,10050,10051,23560,17778,3000,55000,9090,5666,5665,19999,443,8000,8089,6557,8980,9100,9000,8443",
	"firewall": "4444,4433,4443,443,8443",
	"windows":  "88,389,636,593,5985,5986",
	"vpn":      "22,2222,3389,1194,1701,500,4500,1723,5900,5901,5985,5986,443,4443,8443,5938,992,8080,6000,5902",
}

// PortListOrDefault retorna o valor de nmapPS.PortList se não estiver vazio; caso contrário, retorna "".
func (nmapPS *NmapPortScanner) PortListOrDefault() string {
	if nmapPS.PortList != "" {
		return nmapPS.PortList
	}
	return ""
}

// Configure atribui os parâmetros ao scanner.
func (nmapPS *NmapPortScanner) Configure(params PortScanParams) error {
	nmapPS.Targets = params.Targets
	// Define a saída em uma pasta, por exemplo, "portscan"
	nmapPS.OutputFile = filepath.Join(util.PortScanName, params.OutputFile)
	nmapPS.Mode = params.Mode
	nmapPS.Options = params.Options
	nmapPS.PortList = params.PortList
	nmapPS.Category = params.Category
	nmapPS.AllPorts = params.AllPorts
	nmapPS.SimpleScan = params.SimpleScan
	nmapPS.FileMode = params.FileMode

	// Se uma categoria for especificada, mescle-a com a PortList

	if nmapPS.Category != "" {
		nmapPS.PortList = combinePortLists(nmapPS.PortListOrDefault(), nmapPS.Category)

	} else {
		log.Warn().Msgf("Categoria de porta desconhecida: %s", nmapPS.Category)
	}

	// Se AllPorts estiver definido, força o uso de "-p-"
	if nmapPS.AllPorts {
		nmapPS.PortList = ""
	}

	return nil
}

// // mergePortLists combina duas listas de portas, removendo duplicatas.
// func mergePortLists(list1, list2 string) string {
// 	set := make(map[string]bool)
// 	parts1 := strings.Split(list1, ",")
// 	parts2 := strings.Split(list2, ",")
// 	for _, p := range parts1 {
// 		p = strings.TrimSpace(p)
// 		if p != "" {
// 			set[p] = true
// 		}
// 	}
// 	for _, p := range parts2 {
// 		p = strings.TrimSpace(p)
// 		if p != "" {
// 			set[p] = true
// 		}
// 	}
// 	var merged []string
// 	for p := range set {
// 		merged = append(merged, p)
// 	}
// 	return strings.Join(merged, ",")
// }

// buildCommand monta os argumentos para o comando Nmap.
func (nmapPS *NmapPortScanner) buildCommand() (string, []string) {
	outputDir := util.PortScanName
	if err := util.EnsureDir(outputDir); err != nil {
		log.Fatal().Msgf("Error creating directory %s: %v", outputDir, err)
	}

	args := []string{}

	// Escolha de scan: simples (-sS) ou detalhado (-sV -sC)
	if nmapPS.SimpleScan {
		args = append(args, "-sS")
	} else {
		args = append(args, "-sV", "-sC")
	}

	// Adiciona opções extras, se houver.
	if len(nmapPS.Options) > 0 {
		args = append(args, nmapPS.Options...)
	}

	// Adiciona a lista de portas, a menos que AllPorts esteja definido. <=============================== OLHAR
	if !nmapPS.AllPorts {
		if nmapPS.PortList != "" {
			args = append(args, "-p", nmapPS.PortList)
		}
	} else {
		args = append(args, "-p-")
	}

	// Adiciona flag para mostrar somente portas abertas.
	args = append(args, "--open")

	// Adiciona o flag de modo de timing.
	lowerMode := strings.ToLower(strings.TrimSpace(nmapPS.Mode))
	switch lowerMode {
	case "aggressive", "3":
		args = append(args, "-T4")
	case "stealth", "1":
		args = append(args, "-T2")
		// "normal" (2) ou default não adiciona flag adicional.
	}

	// Adiciona os alvos.
	if nmapPS.FileMode {
		// Se os alvos foram passados via arquivo, usa o flag -iL.
		args = append(args, "-iL", nmapPS.Targets[0])
	} else {
		args = append(args, nmapPS.Targets...)
	}

	// Adiciona o comando para gerar os arquivos de saída.
	args = append(args, "-oA", nmapPS.OutputFile)

	commandStr := "nmap " + strings.Join(args, " ")
	return commandStr, args
}

// Execute executa o comando Nmap e retorna a saída bruta.
func (nmapPS *NmapPortScanner) Execute() (string, error) {
	commandStr, args := nmapPS.buildCommand()
	fmt.Printf("%s Running: %s\n", util.MarkerGreen, util.Green(commandStr))

	fmt.Println(args)
	// cmd := exec.Command("nmap", args...)
	// var outputBuffer bytes.Buffer
	// cmd.Stdout = &outputBuffer
	// cmd.Stderr = &outputBuffer

	// if err := cmd.Run(); err != nil {
	// 	return "", fmt.Errorf("nmap execution failed: %w", err)
	// }

	// Lê o arquivo XML gerado pelo Nmap.
	xmlFilePath := nmapPS.OutputFile + ".xml"
	data, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read XML file: %w", err)
	}
	return string(data), nil
}

// Parse processa a saída bruta e extrai informações relevantes (exemplo: lista de portas abertas).
func (nmapPS *NmapPortScanner) Parse(rawOutput string) ([]string, error) {
	// Aqui você pode integrar um parser real (por exemplo, com a biblioteca go-nmap).
	// Este é um exemplo simplificado:
	return []string{"192.168.1.1:80", "192.168.1.2:443"}, nil
}

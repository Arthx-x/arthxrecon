package hostdiscovery

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tomsteele/go-nmap"
)

// MasscanHostDiscovery implementa a interface HostDiscoveryStrategy usando Masscan.
type MasscanHostDiscovery struct {
	Targets []string // Lista de alvos (IPs ou CIDRs) – se FileMode for false, cada elemento é um alvo individual;
	// se FileMode for true, o primeiro elemento deve ser o caminho para o arquivo.
	OutputFile string   // Nome base para os arquivos de saída.
	Mode       string   // Modo do scan (por exemplo, "aggressive", "normal", "passive").
	Options    []string // Outras opções de linha de comando para o Masscan.
	FileMode   bool     // Indica se os targets foram informados via arquivo.
}

// NewMasscanHostDiscovery cria uma instância de MasscanHostDiscovery.
func NewMasscanHostDiscovery() *MasscanHostDiscovery {
	return &MasscanHostDiscovery{}
}

// Configure configura a estratégia com os parâmetros fornecidos.
func (m *MasscanHostDiscovery) Configure(params DiscoveryParams) error {
	// Supondo que DiscoveryParams agora tenha um campo Targets (slice de strings) e FileMode (bool).
	m.Targets = params.Targets
	m.OutputFile = params.OutputFile
	m.Mode = params.Mode
	m.Options = params.Options
	m.FileMode = params.FileMode
	return nil
}

// buildCommand monta o comando masscan com base na configuração.
// Se FileMode for true, utiliza "-iL" com o caminho do arquivo (presente em Targets[0]).
// Caso contrário, adiciona cada target individualmente.
func (m *MasscanHostDiscovery) buildCommand() (string, []string) {
	// Exemplo de comando: masscan -p0-65535 [opções] target(s) -oX outputFile.xml
	args := []string{"-p0-65535"}
	if len(m.Options) > 0 {
		args = append(args, m.Options...)
	}
	if m.FileMode {
		// Se for modo arquivo, o primeiro elemento de Targets é o caminho para o arquivo.
		args = append(args, "-iL", m.Targets[0])
	} else {
		args = append(args, m.Targets...)

	}
	// Masscan gera saída XML com a flag -oX.
	xmlOutput := m.OutputFile + ".xml"
	args = append(args, "-oX", xmlOutput)
	commandStr := "masscan " + strings.Join(args, " ")
	return commandStr, args
}

// Execute executa o comando masscan e, após sua conclusão, lê o arquivo XML gerado e retorna seu conteúdo.
func (m *MasscanHostDiscovery) Execute() (string, error) {
	commandStr, args := m.buildCommand()
	log.Info().Msgf("Executing host discovery with masscan: %s", commandStr)
	cmd := exec.Command("masscan", args...)
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &outputBuffer
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("masscan execution failed: %w", err)
	}
	xmlFilePath := m.OutputFile + ".xml"
	data, err := os.ReadFile(xmlFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read masscan XML file: %w", err)
	}
	return string(data), nil
}

// Parse utiliza a biblioteca go-nmap para converter o XML em estrutura e extrair os hosts.
// Aqui, se o formato do XML gerado pelo masscan for compatível com o do Nmap, podemos reutilizar a mesma lógica.
func (m *MasscanHostDiscovery) Parse(rawOutput string) ([]string, error) {
	result, err := nmap.Parse([]byte(rawOutput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse masscan XML: %w", err)
	}
	var hosts []string
	for _, host := range result.Hosts {
		for _, addr := range host.Addresses {
			if addr.AddrType == "ipv4" {
				hosts = append(hosts, addr.Addr)
				break
			}
		}
	}
	return hosts, nil
}

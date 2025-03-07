package fullrecon

import (
	"fmt"
	"log"

	"github.com/Arthx-x/arthxrecon/internal/hostdiscovery"
	"github.com/Arthx-x/arthxrecon/util"
)

// FullReconOrchestrator coordena a execução dos módulos: host discovery, port scan, enumeration, etc.
type FullReconOrchestrator struct {
	Target string // Pode ser IP, range ou caminho para arquivo com alvos
	Mode   string // Modo do scan (por exemplo, "normal", "aggressive", etc.)
}

// Run executa as etapas de FullRecon em sequência.
func (fr *FullReconOrchestrator) Run() error {
	// Utilize a função comum para obter os targets (a partir de arquivo ou flag única)
	targets := util.GetTargets("", fr.Target)
	if len(targets) == 0 {
		return fmt.Errorf("no valid targets provided")
	}

	// Configura os parâmetros para host discovery.
	params := hostdiscovery.DiscoveryParams{
		Target:     fr.Target,
		OutputFile: "hostdiscovery_output", // Este nome pode ser ajustado ou derivado dinamicamente
		Mode:       fr.Mode,
		Options:    []string{}, // Adicione opções extras se necessário
	}

	// Cria a estratégia: por enquanto, usamos apenas Nmap.
	strategy := hostdiscovery.NewNmapHostDiscovery()

	// Cria o orquestrador de host discovery.
	orchestrator := hostdiscovery.NewHostDiscoveryOrchestrator(strategy, params)

	// Executa o host discovery.
	discoveredHosts, err := orchestrator.Run()
	if err != nil {
		log.Printf("Error during host discovery: %v", err)
		return err
	}

	fmt.Printf("Discovered hosts: %v\n", discoveredHosts)
	// Aqui, você chamaria os outros módulos (portscan, enumeration, vulnanalysis) sequencialmente.
	return nil
}

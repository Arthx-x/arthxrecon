package hostdiscovery

import (
	"fmt"
)

// DiscoveryParams centraliza os parâmetros para a descoberta de hosts.
type DiscoveryParams struct {
	Targets    []string // Lista de alvos (IPs ou CIDRs)
	OutputFile string   // Nome base para os arquivos de saída
	Mode       string   // Modo do scan (aggressive, normal, passive)
	Options    []string // Outras opções, se houver
	FileMode   bool     // Indica se os targets vieram de um arquivo (modo arquivo)
}

// HostDiscoveryStrategy define a interface para uma estratégia de descoberta.
type HostDiscoveryStrategy interface {
	// Configure configura a estratégia com os parâmetros.
	Configure(params DiscoveryParams) error
	// Execute executa a descoberta de hosts e retorna a saída bruta.
	Execute() (string, error)
	// Parse processa a saída bruta e retorna os hosts descobertos.
	Parse(rawOutput string) ([]string, error)
}

// HostDiscoveryOrchestrator coordena a descoberta de hosts usando uma estratégia.
type HostDiscoveryOrchestrator struct {
	Strategy HostDiscoveryStrategy
	Params   DiscoveryParams
}

// NewHostDiscoveryOrchestrator cria um novo orquestrador com a estratégia escolhida.
func NewHostDiscoveryOrchestrator(strategy HostDiscoveryStrategy, params DiscoveryParams) *HostDiscoveryOrchestrator {
	return &HostDiscoveryOrchestrator{
		Strategy: strategy,
		Params:   params,
	}
}

// Run executa o fluxo completo: configura, executa e parseia a saída.
func (orchestrator *HostDiscoveryOrchestrator) Run() ([]string, error) {
	if err := orchestrator.Strategy.Configure(orchestrator.Params); err != nil {
		return nil, fmt.Errorf("failed to configure host discovery: %w", err)
	}

	rawOutput, err := orchestrator.Strategy.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to execute host discovery: %w", err)
	}

	// Se desejar, o parsing pode ser feito em uma goroutine para desacoplar do fluxo principal.
	hosts, err := orchestrator.Strategy.Parse(rawOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host discovery output: %w", err)
	}

	return hosts, nil
}

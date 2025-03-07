package portscan

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// PortScanOrchestrator coordena a execução do port scan usando uma estratégia.
type PortScanOrchestrator struct {
	Strategy PortScanStrategy
	Params   PortScanParams
}

// PortScanStrategy define a interface para uma estratégia de varredura de portas.
type PortScanStrategy interface {
	Configure(params PortScanParams) error
	Execute() (string, error)
	Parse(rawOutput string) ([]string, error)
}

// NewPortScanOrchestrator cria um novo orquestrador com a estratégia escolhida.
func NewPortScanOrchestrator(strategy PortScanStrategy, params PortScanParams) *PortScanOrchestrator {
	return &PortScanOrchestrator{
		Strategy: strategy,
		Params:   params,
	}
}

// Run executa o fluxo completo do port scan: configuração, execução e parsing.
func (orchestrator *PortScanOrchestrator) Run() ([]string, error) {

	if err := orchestrator.Strategy.Configure(orchestrator.Params); err != nil {
		return nil, fmt.Errorf("failed to configure port scan: %w", err)
	}

	rawOutput, err := orchestrator.Strategy.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to execute port scan: %w", err)
	}

	ports, err := orchestrator.Strategy.Parse(rawOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to parse port scan output: %w", err)
	}

	return ports, nil
}

// ===========================================================

// combinePortLists combina a flag --ports com as portas provenientes da flag --category.
// Retorna a união (sem duplicatas) das portas como string, ou "" se nenhum for informado.
func combinePortLists(portList, category string) string {
	set := make(map[int]bool)
	// Processa a flag --ports.
	if portList != "" {
		for _, token := range strings.Split(portList, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if strings.Contains(token, "-") {
				parts := strings.Split(token, "-")
				if len(parts) == 2 {
					start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
					end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
					if err1 == nil && err2 == nil {
						for i := start; i <= end; i++ {
							set[i] = true
						}
					}
				}
			} else {
				if num, err := strconv.Atoi(token); err == nil {
					set[num] = true
				}
			}
		}
	}
	// Processa a flag --category.
	if category != "" {
		catPorts := mergeCategoryPorts(category)
		for _, token := range strings.Split(catPorts, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			if num, err := strconv.Atoi(token); err == nil {
				set[num] = true
			}
		}
	}
	var ports []int
	for p := range set {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	var portStrs []string
	for _, p := range ports {
		portStrs = append(portStrs, fmt.Sprintf("%d", p))
	}
	return strings.Join(portStrs, ",")
}

// mergeCategoryPorts lê as categorias (separadas por vírgula) e retorna uma string com a união de todas as portas,
// sem duplicatas e ordenadas. Se "all" for especificado, inclui todas as portas de todas as categorias.
func mergeCategoryPorts(catStr string) string {
	cats := strings.Split(strings.ToLower(catStr), ",")
	portSet := make(map[int]bool)
	useAll := false
	for _, cat := range cats {
		cat = strings.TrimSpace(cat)
		if cat == "all" {
			useAll = true
			break
		}
	}
	if useAll {
		for _, ports := range portCategories {
			for _, port := range strings.Split(ports, ",") {
				p := strings.TrimSpace(port)
				if p != "" {
					if num, err := strconv.Atoi(p); err == nil {
						portSet[num] = true
					}
				}
			}
		}
	} else {
		for _, cat := range cats {
			cat = strings.TrimSpace(cat)
			if ports, ok := portCategories[cat]; ok {
				for _, port := range strings.Split(ports, ",") {
					p := strings.TrimSpace(port)
					if p != "" {
						if num, err := strconv.Atoi(p); err == nil {
							portSet[num] = true
						}
					}
				}
			} else {
				log.Warn().Msgf("Unknown Category: %s", cat)
			}
		}
	}
	var portSlice []int
	for p := range portSet {
		portSlice = append(portSlice, p)
	}
	sort.Ints(portSlice)
	var portStrs []string
	for _, p := range portSlice {
		portStrs = append(portStrs, fmt.Sprintf("%d", p))
	}
	return strings.Join(portStrs, ",")
}

package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// ParseTargetInput verifica se o valor fornecido na flag -t é um arquivo existente.
// Se for, retorna um slice contendo o próprio caminho e fileMode=true.
// Caso contrário, considera que o valor é uma lista de targets separados por vírgula,
// valida cada um e retorna o slice com fileMode=false.
func ParseTargetInput(input string) (targets []string, fileMode bool) {
	// Verifica se o input corresponde a um arquivo existente.
	if info, err := os.Stat(input); err == nil && !info.IsDir() {
		// É um arquivo – não precisamos ler o conteúdo agora; apenas marcamos fileMode true.
		return []string{input}, true
	}

	// Se não é um arquivo, então trata-se de uma lista de targets.
	parts := strings.Split(input, ",")
	for _, part := range parts {
		t := strings.TrimSpace(part)
		if t == "" {
			continue
		}
		// Valida se o target está no formato esperado (IP ou CIDR).
		if !IsValidTarget(t) {
			log.Fatal().Msgf("%s [%s]", ErrInvalidTarget, t)
		}
		targets = append(targets, t)
	}
	return targets, false
}

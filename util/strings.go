package util

import (
	"regexp"
	"strings"
	"time"
)

// SanitizeString trims spaces from the input string.
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

// SanitizeTarget substitui "/" por "_" para que o target possa ser usado em nomes de arquivos.
func SanitizeTarget(t string) string {
	re := regexp.MustCompile(`[\/]`)
	return re.ReplaceAllString(t, "_")
}

// IsValidTarget valida se o target é um IP válido ou um range em CIDR.
func IsValidTarget(t string) bool {
	regex := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)(\/([0-9]|[12]\d|3[0-2]))?$`)
	return regex.MatchString(t)
}

// SanitizeFileName remove caracteres inválidos de um nome de arquivo.
// Remove caracteres: < > : " / \ | ? *
func SanitizeFileName(name string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	return re.ReplaceAllString(name, "")
}

// GetFormattedTime retorna a data e hora atual formatada no padrão "YYYY/MM/DD HH:MM:SS".
func GetFormattedTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

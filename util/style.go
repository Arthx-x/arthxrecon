package util

import (
	"fmt"

	"github.com/fatih/color"
)

// Variáveis globais para os printers, criadas uma única vez.
var (
	cyanPrinter   = color.New(color.FgCyan).SprintFunc()
	bluePrinter   = color.New(color.FgBlue).SprintFunc()
	greenPrinter  = color.New(color.FgGreen).SprintFunc()
	redPrinter    = color.New(color.FgRed).SprintFunc()
	yellowPrinter = color.New(color.FgYellow).SprintFunc()
	MarkerGreen   = Green("[+]")
	MarkerCyan    = Cyan("[*]")
	MarkerRed     = Red("[-]")
	MarkerYellow  = Yellow("[!]")
)

// Funções que retornam o texto colorido.
func Cyan(text string) string {
	return cyanPrinter(text)
}

func Blue(text string) string {
	return bluePrinter(text)
}

func Green(text string) string {
	return greenPrinter(text)
}

func Red(text string) string {
	return redPrinter(text)
}

func Yellow(text string) string {
	return yellowPrinter(text)
}

// Exemplo de função para exibir um banner com cores
func Banner() {
	// Exemplo usando azul para parte fixa e ciano para partes variáveis.
	// Você pode misturar as funções conforme necessário.
	fmt.Println("")
	fmt.Println("  ┌─────────────────────────────────────────┐ ")
	fmt.Println("  │ █▀█ █▀▄ ▀█▀ █ █ █ █", Cyan("█▀▄ █▀▀ █▀▀ █▀█ █▀█"), "│ ")
	fmt.Println("  │ █▀█ █▀▄  █  █▀█ ▄▀▄", Cyan("█▀▄ █▀▀ █   █ █ █ █"), "│ ")
	fmt.Println("  │ ▀ ▀ ▀ ▀  ▀  ▀ ▀ ▀ ▀", Cyan("▀ ▀ ▀▀▀ ▀▀▀ ▀▀▀ ▀ ▀"), "│ ")
	fmt.Println("  └─────────────────────────────────────────┘ ")
	fmt.Println("\t\t\t\t@Arthx v1.0")
	fmt.Println("")
}

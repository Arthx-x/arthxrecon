package portscan

import (
	"fmt"
	"strings"

	"github.com/Arthx-x/arthxrecon/util"
)

// ShowConfiguration exibe as configurações do port scan de forma centralizada.
func ShowConfiguration(params PortScanParams) {
	fmt.Println("\n┌──────────────────────────────────────────────┐")
	fmt.Printf("  %s: %s\n", util.Green("Target"), strings.Join(params.Targets, ", "))
	fmt.Printf("  %s: %s\n", util.Green("Output"), params.OutputFile)
	fmt.Printf("  %s: %s\n", util.Green("Port Range"), params.PortList)
	fmt.Printf("  %s: %s\n", util.Green("Options"), strings.Join(params.Options, ", "))
	fmt.Printf("  %s: %t\n", util.Green("All Ports"), params.AllPorts)
	fmt.Printf("  %s: %t\n", util.Green("Simple Scan"), params.SimpleScan)
	fmt.Printf("  %s: %s\n", util.Green("Mode"), params.Mode)
	fmt.Printf("  %s: %s\n", util.Green("Category"), params.Category)
	fmt.Println("└──────────────────────────────────────────────┘")
}

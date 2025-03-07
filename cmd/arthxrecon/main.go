package main

import (
	"github.com/Arthx-x/arthxrecon/cmd"
	"github.com/Arthx-x/arthxrecon/util"
)

func main() {
	// Inicializa o logger global, lendo as configurações do config/config.json.
	util.InitializeLogger()
	// Executa o comando raiz (Cobra).
	cmd.Execute()
}

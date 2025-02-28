package main

import (
	"github.com/Arthx-x/arthxrecon/cmd"
	"github.com/Arthx-x/arthxrecon/util/logger"
)

func main() {
	// Inicializa o logger global, lendo as configurações do config/config.json.
	logger.InitializeLogger()
	// Executa o comando raiz (Cobra).
	cmd.Execute()
}

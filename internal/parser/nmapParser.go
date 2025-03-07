package parse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"github.com/tomsteele/go-nmap"
)

// ParseNmapXMLFile reads the Nmap XML file and returns its JSON representation with indentation.
// Any errors are logged using zerolog.
func ParseNmapXMLFile(xmlFilePath string) (string, error) {
	data, err := ioutil.ReadFile(xmlFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read XML file")
		return "", fmt.Errorf("failed to read XML file: %w", err)
	}

	result, err := nmap.Parse(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse XML")
		return "", fmt.Errorf("failed to parse XML: %w", err)
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal JSON")
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	log.Debug().Msg("Successfully parsed Nmap XML and converted to JSON")
	return string(jsonBytes), nil
}

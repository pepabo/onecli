package utils

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

// OutputFormat は出力形式を表します
type OutputFormat string

const (
	OutputFormatYAML OutputFormat = "yaml"
	OutputFormatJSON OutputFormat = "json"
)

// PrintOutput は指定された形式でデータを出力します
func PrintOutput(data interface{}, format OutputFormat) error {
	switch format {
	case OutputFormatJSON:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case OutputFormatYAML:
		fallthrough
	default:
		return yaml.NewEncoder(os.Stdout).Encode(data)
	}
}

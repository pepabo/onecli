package utils

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
)

// OutputFormat は出力形式を表します
type OutputFormat string

const (
	OutputFormatYAML OutputFormat = "yaml"
	OutputFormatJSON OutputFormat = "json"
)

// PrintOutput は指定された形式でデータを出力します
// writerがnilの場合はos.Stdoutを使用します
func PrintOutput(data interface{}, format OutputFormat, writer io.Writer) error {
	if writer == nil {
		writer = os.Stdout
	}

	switch format {
	case OutputFormatJSON:
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		return encoder.Encode(data)
	case OutputFormatYAML:
		return yaml.NewEncoder(writer).Encode(data)
	default:
		return yaml.NewEncoder(writer).Encode(data)
	}
}

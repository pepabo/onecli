package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
)

// OutputFormat は出力形式を表します
type OutputFormat string

const (
	OutputFormatYAML OutputFormat = "yaml"
	OutputFormatJSON OutputFormat = "json"
	OutputFormatCSV  OutputFormat = "csv"
)

// PrintOutput は指定された形式でデータを出力します
// writerがnilの場合はos.Stdoutを使用します
func PrintOutput(data any, format OutputFormat, writer io.Writer) error {
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
	case OutputFormatCSV:
		return encodeCSV(data, writer)
	default:
		return yaml.NewEncoder(writer).Encode(data)
	}
}

// encodeCSV はデータをCSV形式でエンコードします
func encodeCSV(data any, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// データがスライスでない場合はエラー
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return fmt.Errorf("CSV output requires slice data, got %v", val.Kind())
	}

	if val.Len() == 0 {
		return nil
	}

	// 最初の要素からヘッダーを生成
	firstElem := val.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}

	if firstElem.Kind() != reflect.Struct {
		return fmt.Errorf("CSV output requires struct slice, got %v", firstElem.Kind())
	}

	// ヘッダーを生成
	var headers []string
	elemType := firstElem.Type()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		headers = append(headers, field.Name)
	}

	// ヘッダーを書き込み
	if err := csvWriter.Write(headers); err != nil {
		return fmt.Errorf("error writing CSV headers: %v", err)
	}

	// データ行を書き込み
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		var row []string
		for j := 0; j < elem.NumField(); j++ {
			field := elem.Field(j)
			row = append(row, formatFieldValue(field))
		}

		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row %d: %v", i, err)
		}
	}

	return nil
}

// formatFieldValue はフィールドの値を文字列に変換します
func formatFieldValue(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.Struct:
		// time.Timeの場合はRFC3339形式で出力
		if field.Type() == reflect.TypeFor[time.Time]() {
			t := field.Interface().(time.Time)
			return t.Format(time.RFC3339)
		}
		return fmt.Sprintf("%v", field.Interface())
	case reflect.Ptr:
		if field.IsNil() {
			return ""
		}
		return formatFieldValue(field.Elem())
	default:
		return fmt.Sprintf("%v", field.Interface())
	}
}

package headers

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var supportedFormats = map[string]HeaderReader{
	"json": jsonRead,
	"yaml": yamlRead,
	"yml":  yamlRead,
}

func ReadFromFile(path string) (*http.Header, error) {
	format := strings.Replace(filepath.Ext(path), ".", "", -1)
	formatReader, ok := supportedFormats[format]
	if !ok {
		return nil, fmt.Errorf("unsupported format %s", format)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s file:%w", path, err)
	}

	return formatReader(raw)
}

// HeaderReader функция, преобразующие данные из файла в http.Header
type HeaderReader func([]byte) (*http.Header, error)

// SetFormat позволяет выставлять новый метод для преобразования данных из файла в http.Header
// метод перезаписывает существующий функционал для конкретного формата данных
func SetFormat(format string, read HeaderReader) {
	supportedFormats[format] = read
}

// Formats выдаёт все поддерживаемые форматы
func Formats() []string {
	formats := make([]string, 0, len(supportedFormats))
	for f, _ := range supportedFormats {
		formats = append(formats, f)
	}

	return formats
}

func jsonRead(raw []byte) (*http.Header, error) {
	headersMap := make(map[string]string)
	err := json.Unmarshal(raw, &headersMap)
	if err != nil {
		return nil, err
	}
	headers := http.Header{}
	for key, val := range headersMap {
		headers.Set(key, val)
	}

	return &headers, nil
}

func yamlRead(raw []byte) (*http.Header, error) {
	headersMap := make(map[string]string)

	err := yaml.Unmarshal(raw, &headersMap)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}

	for key, value := range headersMap {
		headers.Set(key, value)
	}

	return &headers, nil
}

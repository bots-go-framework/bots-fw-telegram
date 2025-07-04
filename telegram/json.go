package telegram

import (
	"bytes"
	"encoding/json"
)

func getJsonEncoder(buf *bytes.Buffer) *json.Encoder {
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	return encoder
}

func encodeToJsonString(v any) (string, error) {
	var buf bytes.Buffer
	encoder := getJsonEncoder(&buf)
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	return buf.String(), nil
}

package bhformat

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func GoTempate(format string, data any) (string, error) {
	tmpl := template.New("bh").Funcs(sprig.FuncMap())
	tmpl, err := tmpl.Parse(format)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func jsonFormat(v interface{}) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func prettyJsonFormat(v interface{}) (string, error) {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

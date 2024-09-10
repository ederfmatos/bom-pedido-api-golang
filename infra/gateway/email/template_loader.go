package email

import (
	"bytes"
	"html/template"
	"os"
)

type (
	TemplateLoader interface {
		Load(name string, data map[string]string) (string, error)
	}

	htmlTemplateLoader struct {
	}
)

func NewTemplateLoader() TemplateLoader {
	return &htmlTemplateLoader{}
}

func (m *htmlTemplateLoader) Load(name string, data map[string]string) (string, error) {
	templateData, err := os.ReadFile(".resources/templates/" + name + ".html")
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(name).Parse(string(templateData))
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

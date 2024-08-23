package email

import "github.com/cbroglie/mustache"

type (
	TemplateLoader interface {
		Load(name string, data map[string]string) (string, error)
	}

	mustacheLoader struct {
	}
)

func NewTemplateLoader() TemplateLoader {
	return &mustacheLoader{}
}

func (m *mustacheLoader) Load(name string, data map[string]string) (string, error) {
	return mustache.RenderFile(".resources/templates/"+name+".html", data)
}

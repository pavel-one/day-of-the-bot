package templates

import (
	"fmt"

	"github.com/valyala/fasttemplate"
)

// MessageTemplate представляет шаблон сообщения
type MessageTemplate struct {
	template *fasttemplate.Template
}

// TemplateData содержит данные для подстановки в шаблон
type TemplateData map[string]interface{}

// NewTemplate создает новый шаблон из строки
func NewTemplate(templateStr string) (*MessageTemplate, error) {
	template, err := fasttemplate.NewTemplate(templateStr, "{{", "}}")
	if err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return &MessageTemplate{
		template: template,
	}, nil
}

// Execute выполняет шаблон с переданными данными
func (mt *MessageTemplate) Execute(data TemplateData) string {
	return mt.template.ExecuteString(data)
}

// ExecuteString выполняет шаблон напрямую из строки (для простых случаев)
func ExecuteString(templateStr string, data TemplateData) (string, error) {
	template, err := NewTemplate(templateStr)
	if err != nil {
		return "", err
	}

	return template.Execute(data), nil
}

package rendered

import (
	"bytes"
	"github.com/nicolag97/go-mailer/mail"
	"text/template"
)

type TemplateType string

type Template struct {
	Name string
	Type TemplateType
}

const (
	TemplateTypeText       TemplateType = "text"
	TemplateTypeHtml       TemplateType = "html"
	TemplateTypeTextLayout TemplateType = "text_layout"
	TemplateTypeHtmlLayout TemplateType = "html_layout"
)

type RenderedMail struct {
	Sender         mail.Subject
	Recipients     []mail.Subject
	Subject        string
	Attachments    []mail.Attachment
	HtmlLayout     string
	TextLayout     string
	HtmlTemplate   string
	TextTemplate   string
	TemplateFinder func(template Template) []byte
	FuncMap        map[string]interface{}
	ExtraData      map[string]interface{}
}

func (r *RenderedMail) GetSender() mail.Subject {
	return r.Sender
}

func (r *RenderedMail) GetHtmlContent() []byte {
	htmlTmpl := r.TemplateFinder(Template{Name: r.HtmlTemplate, Type: TemplateTypeHtml})
	if len(htmlTmpl) == 0 {
		return []byte("")
	}
	tmpl, err := template.New(string(TemplateTypeHtml)).Funcs(r.FuncMap).Parse(string(htmlTmpl))
	if err != nil {
		return []byte("")
	}
	contentBuffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(contentBuffer, r.ExtraData)
	if err != nil {
		return []byte("")
	}
	htmlLayout := r.TemplateFinder(Template{Name: r.HtmlLayout, Type: TemplateTypeHtmlLayout})
	if len(htmlLayout) == 0 {
		return contentBuffer.Bytes()
	}
	layoutTmpl, err := template.New(string(TemplateTypeHtmlLayout)).Funcs(r.FuncMap).Parse(string(htmlLayout))
	if err != nil {
		return []byte("")
	}
	layoutBuffer := bytes.NewBuffer(nil)
	err = layoutTmpl.Execute(layoutBuffer, struct {
		Content string
	}{Content: contentBuffer.String()})
	if err != nil {
		return []byte("")
	}
	return layoutBuffer.Bytes()
}

func (r *RenderedMail) GetTextContent() []byte {
	textTmpl := r.TemplateFinder(Template{Name: r.TextTemplate, Type: TemplateTypeText})
	if len(textTmpl) == 0 {
		return []byte("")
	}
	tmpl, err := template.New(string(TemplateTypeHtml)).Funcs(r.FuncMap).Parse(string(textTmpl))
	if err != nil {
		return []byte("")
	}
	contentBuffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(contentBuffer, r.ExtraData)
	if err != nil {
		return []byte("")
	}
	textLayout := r.TemplateFinder(Template{Name: r.TextLayout, Type: TemplateTypeTextLayout})
	if len(textLayout) == 0 {
		return contentBuffer.Bytes()
	}
	layoutTmpl, err := template.New(string(TemplateTypeTextLayout)).Funcs(r.FuncMap).Parse(string(textLayout))
	if err != nil {
		return []byte("")
	}
	layoutBuffer := bytes.NewBuffer(nil)
	err = layoutTmpl.Execute(layoutBuffer, struct {
		Content string
	}{Content: contentBuffer.String()})
	if err != nil {
		return []byte("")
	}
	return layoutBuffer.Bytes()
}

func (r *RenderedMail) GetRecipients() []mail.Subject {
	return r.Recipients
}

func (r *RenderedMail) GetSubject() string {
	return r.Subject
}

func (r *RenderedMail) GetAttachments() []mail.Attachment {
	return r.Attachments
}

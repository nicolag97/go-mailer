package rendered

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nicolag97/go-mailer/mail"
	"github.com/opentracing/opentracing-go"
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
	Layout         string
	Template       string
	TemplateFinder func(template Template) []byte
	FuncMap        map[string]interface{}
	ExtraData      map[string]interface{}
}

func (r *RenderedMail) GetSender() mail.Subject {
	return r.Sender
}

func (r *RenderedMail) GetHtmlContent(ctx context.Context) []byte {
	var newSpan opentracing.Span
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		tracer := parent.Tracer()
		newSpan = tracer.StartSpan("Mail.RenderHtml", opentracing.ChildOf(pctx))
		defer newSpan.Finish()
	}
	htmlTmpl := r.TemplateFinder(Template{Name: r.Template, Type: TemplateTypeHtml})
	if len(htmlTmpl) == 0 {
		newSpan.LogEvent("[ERROR]: No body")
		return []byte("")
	}
	tmpl, err := template.New(string(TemplateTypeHtml)).Funcs(r.FuncMap).Parse(string(htmlTmpl))
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	contentBuffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(contentBuffer, r.ExtraData)
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	htmlLayout := r.TemplateFinder(Template{Name: r.Layout, Type: TemplateTypeHtmlLayout})
	if len(htmlLayout) == 0 {
		return contentBuffer.Bytes()
	}
	layoutTmpl, err := template.New(string(TemplateTypeHtmlLayout)).Funcs(r.FuncMap).Parse(string(htmlLayout))
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	layoutBuffer := bytes.NewBuffer(nil)
	err = layoutTmpl.Execute(layoutBuffer, struct {
		Content string
	}{Content: contentBuffer.String()})
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	return layoutBuffer.Bytes()
}

func (r *RenderedMail) GetTextContent(ctx context.Context) []byte {
	var newSpan opentracing.Span
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		tracer := parent.Tracer()
		newSpan = tracer.StartSpan("Mail.RenderPlainText", opentracing.ChildOf(pctx))
		defer newSpan.Finish()
	}
	textTmpl := r.TemplateFinder(Template{Name: r.Template, Type: TemplateTypeText})
	if len(textTmpl) == 0 {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: No template"))
		return []byte("")
	}
	tmpl, err := template.New(string(TemplateTypeHtml)).Funcs(r.FuncMap).Parse(string(textTmpl))
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	contentBuffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(contentBuffer, r.ExtraData)
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	textLayout := r.TemplateFinder(Template{Name: r.Layout, Type: TemplateTypeTextLayout})
	if len(textLayout) == 0 {
		return contentBuffer.Bytes()
	}
	layoutTmpl, err := template.New(string(TemplateTypeTextLayout)).Funcs(r.FuncMap).Parse(string(textLayout))
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
		return []byte("")
	}
	layoutBuffer := bytes.NewBuffer(nil)
	err = layoutTmpl.Execute(layoutBuffer, struct {
		Content string
	}{Content: contentBuffer.String()})
	if err != nil {
		newSpan.LogEvent(fmt.Sprintf("[ERROR]: %v", err.Error()))
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

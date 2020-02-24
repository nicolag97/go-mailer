package examples

import (
	"fmt"
	"github.com/nicolag97/go-mailer/mail"
	"github.com/nicolag97/go-mailer/mail/rendered"
	"github.com/nicolag97/go-mailer/mailer/smtp"
	"io/ioutil"
	"log"
)

func TemplateFinder(template rendered.Template) []byte {
	switch template.Type {
	case rendered.TemplateTypeHtml:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/html/%v", template.Name))
		if err != nil {
			return []byte("")
		}
		return file
	case rendered.TemplateTypeText:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/text/%v", template.Name))
		if err != nil {
			return []byte("")
		}
		return file
	case rendered.TemplateTypeHtmlLayout, rendered.TemplateTypeTextLayout:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/layout/%v", template.Name))
		if err != nil {
			return []byte("")
		}
		return file
	default:
		return []byte("")
	}
}

func SendRenderedEmailSmtp() {
	Client, err := smtp.NewSmtpClient("smtp.example.com", "465", smtp.SmtpCredentials{
		Username: "test@example.com",
		Password: "password",
		Identity: "test@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	mail := rendered.RenderedMail{
		Sender: mail.Subject{
			Name: "Example",
			Mail: "test@example.com",
		},
		Recipients: []mail.Subject{
			{
				Name: "TestFoo",
				Mail: "foo@example.com",
			},
			{
				Name: "TestBar",
				Mail: "bar@example.com",
			},
		},
		Subject:        "Hi, I'm a test email.",
		HtmlLayout:     "layout.html.tmpl",
		TextLayout:     "layout.text.tmpl",
		HtmlTemplate:   "html.tmpl",
		TextTemplate:   "text.tmpl",
		TemplateFinder: TemplateFinder,
		FuncMap:        nil,
		ExtraData:      nil,
	}
	err = Client.Send(&mail)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Mail delivered Successful")
}

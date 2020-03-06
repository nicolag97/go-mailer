package examples

import (
	"fmt"
	"github.com/nicolag97/go-mailer"
	"github.com/nicolag97/go-mailer/mail"
	"github.com/nicolag97/go-mailer/mail/rendered"
	"io/ioutil"
	"log"
)

func TemplateFinder(template rendered.Template) []byte {
	switch template.Type {
	case rendered.TemplateTypeHtml:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/html/%v", fmt.Sprintf(template.Name, "html")))
		if err != nil {
			return []byte("")
		}
		return file
	case rendered.TemplateTypeText:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/text/%v", fmt.Sprintf(template.Name, "text")))
		if err != nil {
			return []byte("")
		}
		return file
	case rendered.TemplateTypeHtmlLayout:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/layout/%v", fmt.Sprintf(template.Name, "html")))
		if err != nil {
			return []byte("")
		}
		return file
	case rendered.TemplateTypeTextLayout:
		file, err := ioutil.ReadFile(fmt.Sprintf("templates/layout/%v", fmt.Sprintf(template.Name, "text")))
		if err != nil {
			return []byte("")
		}
		return file
	default:
		return []byte("")
	}
}

func SendRenderedEmailSmtp() {
	Client, err := go_mailer.NewSmtpClient("smtp.example.com", "465", go_mailer.SmtpCredentials{
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
		Layout:         "layout.%v.tmpl",
		Template:       "%v.tmpl",
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

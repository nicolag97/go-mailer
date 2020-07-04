package examples

import (
	"fmt"
	"github.com/nicolag97/go-mailer"
	"github.com/nicolag97/go-mailer/mail"
	"github.com/nicolag97/go-mailer/mail/simple"
	"log"
)

func SendSampleMailSmtp() {
	Client, err := go_mailer.NewSmtpClient("smtp.example.com", "465", go_mailer.SmtpCredentials{
		Username: "test@example.com",
		Password: "password",
		Identity: "test@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}
	Mail := simple.SimpleMail{
		Sender: mail.Subject{
			Name: "Example",
			Mail: "test@example.com",
		},
		Html: []byte("<b>Hi, I'm a test Mail in HTML.</b>"),
		Text: []byte("Hi, I'm a test Mail in plain text."),
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
		Subject: "Hi, I'm a a test Email.",
	}
	err = Client.Send(&Mail)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Mail delivered Successful")
}

func SendSampleMailSES() {
	Client, err := go_mailer.NewSesMailer("TestUsername", "TestPassword", "TestRegion")
	if err != nil {
		log.Fatal(err)
	}
	mail := simple.SimpleMail{
		Sender: mail.Subject{
			Name: "Example",
			Mail: "test@example.com",
		},
		Html: []byte("<b>Hi, I'm a test mail in HTML.</b>"),
		Text: []byte("Hi, I'm a test mail in plain text."),
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
		Subject: "Hi, I'm a a test Email.",
	}
	err = Client.Send(&mail)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Mail delivered Successful")
}

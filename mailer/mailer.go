package mailer

import "github.com/nicolag97/go-mailer/mail"

type MailClient interface {
	Send(mail mail.Mail) error
}

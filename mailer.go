package go_mailer

import (
	"context"
	"github.com/nicolag97/go-mailer/mail"
)

type MailClient interface {
	Send(ctx context.Context, Mail mail.Mail) error
}

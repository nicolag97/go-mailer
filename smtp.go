package go_mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/labstack/gommon/random"
	"github.com/nicolag97/go-mailer/mail"
	"log"
	"net/smtp"
)

type SmtpMailer struct {
	auth   smtp.Auth
	conn   SmtpConnection
	client *smtp.Client
}

func (s *SmtpMailer) Send(ctx context.Context, mail mail.Mail) error {
	textContent := mail.GetTextContent(context.Background())
	htmlContent := mail.GetHtmlContent(context.Background())
	if (len(textContent) == 0) && (len(htmlContent) == 0) {
		return errors.New("No body provided")
	}

	MailCtx := RawMailContext{
		To:                  mail.GetRecipients(),
		From:                mail.GetSender(),
		Subject:             mail.GetSubject(),
		MixedBoundary:       random.String(10),
		AlternativeBoundary: random.String(10),
		Parts: []RawMsgPart{{ContentType: MimeTypePlain, Message: string(textContent)},
			{ContentType: MimeTypeHtml, Message: string(htmlContent)}},
		Attachments: GetRawMailAttachments(mail.GetAttachments()),
	}
	content, err := RenderRawMail(MailCtx)
	if err != nil {
		return err
	}
	defer s.client.Quit()
	err = s.client.Mail(mail.GetSender().Mail)
	if err != nil {
		return err
	}
	for _, v := range mail.GetRecipients() {
		err = s.client.Rcpt(v.Mail)
		if err != nil {
			return err
		}
	}
	w, err := s.client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(content)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewSmtpClient(Host string, Port string, Creds SmtpCredentials) (*SmtpMailer, error) {
	Auth := smtp.PlainAuth("", Creds.Username, Creds.Password, Host)
	connParams := SmtpConnection{
		Host: Host,
		Port: Port,
	}
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         Host,
	}
	conn, err := tls.Dial("tcp", connParams.GetAddr(), tlsconfig)
	if err != nil {
		log.Panic(err)
	}
	c, err := smtp.NewClient(conn, connParams.Host)
	if err != nil {
		return nil, err
	}
	err = c.Auth(Auth)
	if err != nil {
		return nil, err
	}
	return &SmtpMailer{
		auth:   Auth,
		conn:   connParams,
		client: c,
	}, nil
}

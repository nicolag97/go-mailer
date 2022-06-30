package go_mailer

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/labstack/gommon/random"
	"github.com/nicolag97/go-mailer/mail"
)

type SesMailer struct {
	AwsSes *session.Session
	AwsSVC *ses.SES
}

func (s *SesMailer) Send(ctx context.Context, mail mail.Mail) error {
	if len(mail.GetAttachments()) > 0 {
		ctx := RawMailContext{
			To:                  mail.GetRecipients(),
			From:                mail.GetSender(),
			Subject:             mail.GetSubject(),
			AlternativeBoundary: random.String(10),
			MixedBoundary:       random.String(10),
			Parts: []RawMsgPart{{ContentType: MimeTypePlain, Message: string(mail.GetTextContent(context.Background()))},
				{ContentType: MimeTypeHtml, Message: string(mail.GetHtmlContent(context.Background()))}},
			Attachments: GetRawMailAttachments(mail.GetAttachments()),
		}
		content, err := RenderRawMail(ctx)
		if err != nil {
			return err
		}
		input := &ses.SendRawEmailInput{
			Destinations: GetToAddresses(mail.GetRecipients()),
			RawMessage:   &ses.RawMessage{Data: content},
			Source:       aws.String(mail.GetSender().Mail),
		}
		err = input.Validate()
		_, err = s.AwsSVC.SendRawEmail(input)
		if err != nil {
			return err
		}
		return nil
	}
	textContent := mail.GetTextContent(context.Background())
	htmlContent := mail.GetHtmlContent(context.Background())
	if (len(textContent) == 0) && (len(htmlContent) == 0) {
		return errors.New("No body provided")
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: GetToAddresses(mail.GetRecipients()),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(defaultCharSet),
					Data:    aws.String(string(htmlContent)),
				},
				Text: &ses.Content{
					Charset: aws.String(defaultCharSet),
					Data:    aws.String(string(textContent)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(defaultCharSet),
				Data:    aws.String(mail.GetSubject()),
			},
		},
		Source: aws.String(mail.GetSender().Mail),
	}
	_, err := s.AwsSVC.SendEmail(input)
	if err != nil {
		return err
	}
	return nil
}

func NewSesMailer(UserName string, Password string, Region string) (*SesMailer, error) {
	Credentials := credentials.NewStaticCredentials(UserName, Password, "")
	Session, err := session.NewSession(&aws.Config{
		Credentials: Credentials,
		Region:      aws.String(Region),
	})
	if err != nil {
		return nil, err
	}
	Svc := ses.New(Session)
	return &SesMailer{
		AwsSes: Session,
		AwsSVC: Svc,
	}, nil
}

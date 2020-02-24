package simple

import "github.com/nicolag97/go-mailer/mail"

type SimpleMail struct {
	Sender      mail.Subject
	Html        []byte
	Text        []byte
	Recipients  []mail.Subject
	Subject     string
	Attachments []mail.Attachment
}

func (s *SimpleMail) GetSender() mail.Subject {
	return s.Sender
}

func (s *SimpleMail) GetHtmlContent() []byte {
	return s.Html
}

func (s *SimpleMail) GetTextContent() []byte {
	return s.Text
}

func (s *SimpleMail) GetRecipients() []mail.Subject {
	return s.Recipients
}

func (s *SimpleMail) GetSubject() string {
	return s.Subject
}

func (s *SimpleMail) GetAttachments() []mail.Attachment {
	return s.Attachments
}

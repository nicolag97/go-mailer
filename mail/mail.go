package mail

import "context"

type Subject struct {
	Name string
	Mail string
}

type Attachment struct {
	Name        string
	Content     []byte
	ContentType string
}

type Mail interface {
	GetSender() Subject
	GetHtmlContent(ctx context.Context) []byte
	GetTextContent(ctx context.Context) []byte
	GetRecipients() []Subject
	GetSubject() string
	GetAttachments() []Attachment
}

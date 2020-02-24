package mail

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
	GetHtmlContent() []byte
	GetTextContent() []byte
	GetRecipients() []Subject
	GetSubject() string
	GetAttachments() []Attachment
}

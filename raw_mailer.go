package go_mailer

import (
	"bytes"
	"encoding/base64"
	"github.com/nicolag97/go-mailer/mail"
	"text/template"
)

const RawMailTmpl = `From: {{.From.Name}} <{{.From.Mail}}> 
To:  {{range $to:=.To}}{{$to.Name}} <{{$to.Mail}}> {{end}}
Subject: {{.Subject}}
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="{{.MixedBoundary}}"

--{{.MixedBoundary}}
Content-type: multipart/alternative; boundary="{{.AlternativeBoundary}}"
{{$alt_boundary:=.AlternativeBoundary}}{{range $part:= .Parts}}
--{{$alt_boundary}}
Content-type: {{$part.ContentType}}; charset="UTF-8"

{{$part.Message}}
{{end}}{{$mix_boundary := .MixedBoundary}}
--{{.AlternativeBoundary}}--
{{range $attach:= .Attachments}}
--{{$mix_boundary}}
Content-Type: {{$attach.ContentType}};name="{{$attach.Name}}"
Content-Transfer-Encoding: base64

{{$attach.Message}}
{{end}}
--{{.MixedBoundary}}--
`

func RenderRawMail(ctx RawMailContext) ([]byte, error) {
	tmpl, err := template.New("mail").Parse(RawMailTmpl)
	if err != nil {
		return []byte(""), err
	}
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, ctx)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

type RawMsgPart struct {
	ContentType string
	Message     string
}

type RawMsgAttachments struct {
	Name        string
	ContentType string
	Message     string
}

type RawMailContext struct {
	To                  []mail.Subject
	From                mail.Subject
	Subject             string
	AlternativeBoundary string
	MixedBoundary       string
	Parts               []RawMsgPart
	Attachments         []RawMsgAttachments
}

func GetRawMailAttachments(attachment []mail.Attachment) []RawMsgAttachments {
	var result []RawMsgAttachments
	for _, v := range attachment {
		result = append(result, RawMsgAttachments{
			Name:        v.Name,
			ContentType: v.ContentType,
			Message:     base64.StdEncoding.EncodeToString(v.Content),
		})
	}
	return result
}

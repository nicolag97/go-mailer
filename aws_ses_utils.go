package go_mailer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/nicolag97/go-mailer/mail"
)

const (
	defaultCharSet = "UTF-8"
)

func GetToAddresses(sub []mail.Subject) []*string {
	var toAddrs []*string
	for _, v := range sub {
		toAddrs = append(toAddrs, aws.String(v.Mail))
	}
	return toAddrs
}

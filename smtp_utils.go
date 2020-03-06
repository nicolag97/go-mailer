package go_mailer

import (
	"fmt"
)

type SmtpCredentials struct {
	Username string
	Password string
	Identity string
}

type SmtpConnection struct {
	Host string
	Port string
}

func (s *SmtpConnection) GetAddr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

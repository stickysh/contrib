/*
Copyright 2015 Sticky Contrib Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package basicemail

import "net/smtp"

type conn struct {
	mime    string
	address string
	auth    smtp.Auth
}

type Requester interface {
	Subject() string
	Body() string
	To() []string
	From() string
}

// NewSender create
func NewSender(mime string, smtpAddress string, authentication smtp.Auth) *conn {
	return &conn{
		mime:    mime,
		address: smtpAddress,
		auth:    authentication,
	}
}

func (c *conn) SendEmail(to []string, from, subject, body string) error {
	subject = "Subject: " + subject + "!\n"
	msg := []byte(subject + c.mime + "\n" + body)
	return smtp.SendMail(c.address, c.auth, from, to, msg)
}
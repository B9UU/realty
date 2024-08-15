package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer  *mail.Dialer
	sender  string
	Message struct {
		appName string
		Method  string
		URI     string
	}
}

func New(host, sender, username, password string, port int) Mailer {
	dailer := mail.NewDialer(host, port, username, password)
	// set new timeout
	dailer.Timeout = time.Second * 5
	return Mailer{
		dialer: dailer,
		sender: sender,
	}
}
func (m Mailer) Send(recipient, templateFile string, data interface{}) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// try to send the email 3 times
	for i := 0; i < 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// if everything worked, return nil
		if nil == err {
			return nil
		}
		// else sleep for 500ms
		time.Sleep(time.Millisecond * 500)
	}
	return err
}

package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"donetick.com/core/config"
	"donetick.com/core/internal/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	gomail "gopkg.in/gomail.v2"
)

type EmailSender struct {
	client  *gomail.Dialer
	appHost string
	bundle  *i18n.Bundle
}

func NewEmailSender(conf *config.Config, bundle *i18n.Bundle) *EmailSender {
	client := gomail.NewDialer(conf.EmailConfig.Host, conf.EmailConfig.Port, conf.EmailConfig.Email, conf.EmailConfig.Key)
	return &EmailSender{
		client:  client,
		appHost: conf.EmailConfig.AppHost,
		bundle:  bundle,
	}
}

func (es *EmailSender) getLocalizer(c context.Context) *i18n.Localizer {
	localizer, ok := c.Value(middleware.I18nKey).(*i18n.Localizer)
	if !ok {
		// fallback to english
		return i18n.NewLocalizer(es.bundle, "en")
	}
	return localizer
}

func (es *EmailSender) SendVerificationEmail(c context.Context, to, code string) error {
	localizer := es.getLocalizer(c)
	subject, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "welcome_email_subject"})
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", es.client.Username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	lang := localizer.LanguageTags()[0].String()
	templatePath := filepath.Join("locales", lang, "verification_email.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// fallback to english
		templatePath = filepath.Join("locales", "en", "verification_email.html")
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	u := es.appHost + "/verify?c=" + encodeEmailAndCode(to, code)
	data := struct {
		VerifyURL string
	}{
		VerifyURL: u,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	msg.SetBody("text/html", body.String())

	err = es.client.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func (es *EmailSender) SendResetPasswordEmail(c context.Context, to, code string) error {
	localizer := es.getLocalizer(c)
	subject, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: "reset_password_email_subject"})
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", es.client.Username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	lang := localizer.LanguageTags()[0].String()
	templatePath := filepath.Join("locales", lang, "password_reset_email.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		// fallback to english
		templatePath = filepath.Join("locales", "en", "password_reset_email.html")
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	u := es.appHost + "/password/update?c=" + encodeEmailAndCode(to, code)
	data := struct {
		VerifyURL string
	}{
		VerifyURL: u,
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	msg.SetBody("text/html", body.String())

	err = es.client.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func encodeEmailAndCode(email, code string) string {
	data := email + ":" + code
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeEmailAndCode(encoded string) (string, string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}
	parts := string(data)
	split := strings.Split(parts, ":")
	if len(split) != 2 {
		return "", "", fmt.Errorf("Invalid format")
	}
	return split[0], split[1], nil
}

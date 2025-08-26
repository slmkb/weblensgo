package models

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type EmailService struct {
	SMTPServer *mail.SMTPServer
}

func NewSMTPClient(smtpConfig SMTPConfig) (*mail.SMTPServer, error) {
	server := mail.NewSMTPClient()
	server.Host = smtpConfig.Host
	p, err := strconv.Atoi(smtpConfig.Port)
	if err != nil {
		return nil, fmt.Errorf("es.Connect: %w", err)
	}
	server.Port = p
	server.Username = smtpConfig.User
	server.Password = smtpConfig.Password
	server.Encryption = mail.EncryptionSTARTTLS

	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return server, nil
}

func (es *EmailService) ForgotPassword(pr *PasswordReset) error {
	smtpClient, err := es.SMTPServer.Connect()

	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()
	email.SetFrom("Support WebLensGo <support@example.com>")
	email.AddTo(pr.Email)
	email.SetSubject("Password reset link")
	email.SetBody(mail.TextHTML, fmt.Sprintf(`<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Password reset link</title>
	</head>
	<body>
		<p>To reset your password please visit the following link.</p>
		<p><a href="%s">http://eaasd.com/reset-password?token=%[1]s</a></p>
	</body>
</html>`, pr.RToken))
	email.AddAlternative(
		mail.TextPlain,
		fmt.Sprintf("To reset your password please visit the following link.\n http://exsle.com/update-password?token=%s", pr.RToken),
	)
	email.SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

	if email.Error != nil {
		return fmt.Errorf("es.Send: %w", err)
	}

	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("es.Send: %w", err)
	} else {
		log.Println("Email Sent")
	}
	return nil
}

func (es *EmailService) Send(email mail.Email) error {
	// // you can add dkim signature to the email.
	// // to add dkim, you need a private key already created one.
	// if privateKey != "" {
	// 	options := dkim.NewSigOptions()
	// 	options.PrivateKey = []byte(privateKey)
	// 	options.Domain = "example.com"
	// 	options.Selector = "default"
	// 	options.SignatureExpireIn = 3600
	// 	options.Headers = []string{"from", "date", "mime-version", "received", "received"}
	// 	options.AddSignatureTimestamp = true
	// 	options.Canonicalization = "relaxed/relaxed"

	// 	email.SetDkim(options)
	// }

	// always check error after send
	// if email.Error != nil {
	// 	return fmt.Errorf("es.Send: %w", err)
	// }

	// // Call Send and pass the client
	// err = email.Send(smtpClient)
	// if err != nil {
	// 	return fmt.Errorf("es.Send: %w", err)
	// } else {
	// 	log.Println("Email Sent")
	// }
	return nil
}

package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type SmtpConfig struct {
	To             string
	From           string
	Port           string
	Server         string
	Subject        string
	Body           string
	Username       string
	Password       string
	AttachmentPath string
}

// Submit the turn by shelling out to mailsend
func (mailConfig SmtpConfig) SubmitTurnMailsend() error {
	// export SMTP_USER_PASS=password
	cmdName := "mailsend"
	err := exec.Command(cmdName, "-V").Run()
	if err != nil {
		return errors.New("Could not find 'mailsend'. Please see TODO")
	}

	cmdArgs := []string{"-to", mailConfig.To, "-from", mailConfig.From, "-starttls", "-port", mailConfig.Port, "-auth", "-smtp", mailConfig.Server, "-sub", mailConfig.Subject, "+cc", "+bc", "-user", mailConfig.Username, "-M", " ", "-mime-type", "application/octet-stream", "-attach", mailConfig.AttachmentPath}

	cmd := exec.Command(cmdName, cmdArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send mail: %v", err.Error()))
	}

	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-time.After(10 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return errors.New(fmt.Sprintf("Timeout while sending mail: %s", err.Error()))
		}
		<-cmdDone
	case err = <-cmdDone:
		if err != nil {
			return errors.New(fmt.Sprintf("Failed sending mail: %s", err.Error()))
		}
	}

	return nil
}

// Submit the turn by using a builtin mailer
func (mailConfig SmtpConfig) SubmitTurnBuiltin() error {
	message := gomail.NewMessage()
	message.SetHeader("From", mailConfig.From)
	message.SetHeader("To", mailConfig.To)
	message.SetHeader("Subject", mailConfig.Subject)
	message.SetBody("text/plain", mailConfig.Body)
	message.Attach(mailConfig.AttachmentPath)

	port, err := strconv.Atoi(mailConfig.Port)
	if err != nil {
		return err
	}

	dialer := gomail.NewPlainDialer(mailConfig.Server, port, mailConfig.Username, mailConfig.Password)

	if err := dialer.DialAndSend(message); err != nil {
		return (err)
	}

	return nil
}

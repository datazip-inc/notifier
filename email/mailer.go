package email

func NotifyEmail(subject string, fromEmail string, recipient Recipient, content string) {
	standardMailer.NotifyEmail(subject, fromEmail, recipient, content)
}

func NotifyEmailE(subject string, fromEmail string, recipient Recipient, content string) error {
	return standardMailer.NotifyEmailE(subject, fromEmail, recipient, content)
}

package notifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Notify error logs on Slack
func (n *Notifier) NotifyError(errorAt, description, errString string, fields ...string) {
	if err := n.NotifyErrorE(errorAt, description, errString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🔴 Error log reported on slack at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifySuccess(successAt, description, successString string, fields ...string) {
	if err := n.NotifySuccessE(successAt, description, successString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟢 Success log reported on slack at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifyWarn(warnAt, description, warnString string, fields ...string) {
	if err := n.NotifyWarnE(warnAt, description, warnString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟡 Warn log reported on slack at %s", time.Now())
}

// Notify error logs on Slack and returns error
func (n *Notifier) NotifyErrorE(errorAt, description, errString string, fields ...string) error {
	if !isConfigured(n.config.Default) {
		return errors.New("❌ Slack default config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	mainMessageFields := []slack.AttachmentField{
		{Title: "ErrorAt", Value: errorAt},
		{Title: "Description", Value: description},
	}

	additionalMessageFields := []slack.AttachmentField{}

	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	if len(errString) < 4000 {
		err := n.SendOnSlack(errString, errorColor, n.config.Default, mainMessageFields, additionalMessageFields)
		if err != nil {
			return fmt.Errorf("❌ Failed to report error on slack: %s", err)
		}
	} else {
		err := n.SendOnSlackAsFile(errString, errorColor, n.config.Default, mainMessageFields, additionalMessageFields)
		if err != nil {
			return fmt.Errorf("❌ Failed to report error on slack: %s", err)
		}
	}

	return nil
}

// Notify success logs on Slack and returns error
func (n *Notifier) NotifySuccessE(successAt, description, successString string, fields ...string) error {
	if !isConfigured(n.config.Default) {
		return errors.New("❌ Slack default config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	mainMessageFields := []slack.AttachmentField{
		{Title: "SuccessAt", Value: successAt},
		{Title: "Description", Value: description},
	}

	additionalMessageFields := []slack.AttachmentField{}
	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	err := n.SendOnSlack(successString, successColor, n.config.Default, mainMessageFields, additionalMessageFields)
	if err != nil {
		return fmt.Errorf("❌ Failed to report success on slack: %s", err)
	}

	return nil
}

// Notify warn logs on Slack
func (n *Notifier) NotifyWarnE(warnAt, description, warnString string, fields ...string) error {
	if !isConfigured(n.config.Debug) {
		return errors.New("❌ Slack debug config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	mainMessageFields := []slack.AttachmentField{
		{Title: "WarnAt", Value: warnAt},
		{Title: "Description", Value: description},
	}

	additionalMessageFields := []slack.AttachmentField{}

	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, slack.AttachmentField{
			Title: fields[i],
			Value: fields[i+1],
		})
	}

	err := n.SendOnSlack(warnString, warnColor, n.config.Debug, mainMessageFields, additionalMessageFields)
	if err != nil {
		return fmt.Errorf("❌ Failed to report warn on slack: %s", err)
	}

	return nil
}

// SendOnSlackAsFile sends text as file on slack channel
func (n *Notifier) SendOnSlackAsFile(text, messageColor string, channelConfig *SlackChannelConfig, messageAttachments []slack.AttachmentField, additionalMessageFields []slack.AttachmentField) error {
	err := n.SendOnSlack("", messageColor, channelConfig, messageAttachments, additionalMessageFields)
	if err != nil {
		return err
	}

	// Create the Slack attachment that we will send to the channel
	fileattachment := slack.FileUploadParameters{
		Content:  text,
		Channels: []string{channelConfig.ChannelID},
	}

	_, err = n.slackClient.UploadFile(fileattachment)
	if err != nil {
		return fmt.Errorf("failed to upload file : %s", err)
	}

	return nil
}

func (n *Notifier) SendOnSlack(text, messageColor string, channelConfig *SlackChannelConfig, messageAttachments []slack.AttachmentField, additionalMessageFields []slack.AttachmentField) error {
	mentions := generateMentions(channelConfig.Mentions)

	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Pretext: mentions,
		Text:    text,
		Color:   messageColor,
		Fields:  messageAttachments,
		Footer:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_, timeStamp, err := n.slackClient.PostMessage(
		channelConfig.ChannelID,
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return fmt.Errorf("failed to send message : %s", err)
	}

	for _, field := range additionalMessageFields {
		newAttachment := slack.Attachment{
			Color:  messageColor,
			Fields: []slack.AttachmentField{field},
			Footer: time.Now().Format("2006-01-02 15:04:05"),
		}
		_, currentTimestamp, err := n.slackClient.PostMessage(
			channelConfig.ChannelID,
			slack.MsgOptionAttachments(newAttachment),
			slack.MsgOptionTS(timeStamp),
		)
		if err != nil {
			continue
		}

		// update the timestamp if message sent in thread successfully
		timeStamp = currentTimestamp
	}

	return nil
}

func isConfigured(channelConfig *SlackChannelConfig) bool {
	if channelConfig == nil {
		return false
	} else if channelConfig.ChannelID == "" {
		return false
	}

	return true
}

func generateMentions(users []string) string {
	str := ""
	for _, user := range users {
		str += fmt.Sprintf("<@%s> ", user)
	}

	return str
}

package notifier

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

// Notify error logs on Slack
func (n *Notifier) NotifyError(errorAt, description, errString string, fields ...string) {
	if err := n.NotifyErrorE(errorAt, description, errString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🔴 Error log reported on discord at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifySuccess(successAt, description, successString string, fields ...string) {
	if err := n.NotifySuccessE(successAt, description, successString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟢 Success log reported on discord at %s", time.Now())
}

// Notify success logs on Slack
func (n *Notifier) NotifyWarn(warnAt, description, warnString string, fields ...string) {
	if err := n.NotifyWarnE(warnAt, description, warnString, fields...); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("🟡 Warn log reported on discord at %s", time.Now())
}

// notify error logs on discord
func (n *Notifier) NotifyErrorE(errorAt, description, errString string, fields ...string) error {
	if !isConfigured(n.config.Error) {
		return errors.New("❌ Discord error config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	messageEmbeds := []*discordgo.MessageEmbedField{
		{Name: "ErrorAt", Value: errorAt},
		{Name: "Description", Value: description},
	}

	additionalMessageFields := []*discordgo.MessageEmbed{}

	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, &discordgo.MessageEmbed{
			Title:       fields[i],
			Description: fields[i+1],
		})
	}

	err := n.SendOnDiscord(errString, errorColorINT, n.config.Error, messageEmbeds, additionalMessageFields)
	if err != nil {
		return fmt.Errorf("❌ Failed to report error on Discord: %s", err)
	}

	return nil
}

// notify success logs on discord
func (n *Notifier) NotifySuccessE(successAt, description, successString string, fields ...string) error {
	if !isConfigured(n.config.Success) {
		return errors.New("❌ Discord success config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	messageEmbeds := []*discordgo.MessageEmbedField{
		{Name: "SuccessAt", Value: successAt},
		{Name: "Description", Value: description},
	}

	additionalMessageFields := []*discordgo.MessageEmbed{}

	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, &discordgo.MessageEmbed{
			Title:       fields[i],
			Description: fields[i+1],
		})
	}

	err := n.SendOnDiscord(successString, successColorINT, n.config.Success, messageEmbeds, additionalMessageFields)
	if err != nil {
		return fmt.Errorf("❌ Failed to report success on Discord: %s", err)
	}

	return nil
}

// notify warn logs on discord
func (n *Notifier) NotifyWarnE(warnAt, description, warnString string, fields ...string) error {
	if !isConfigured(n.config.Warn) {
		return errors.New("❌ Discord warn config not found or not properly configured")
	}

	if len(fields)%2 != 0 {
		return errors.New("❌ Invalid number of fields passed, only even number of fields allowed")
	}

	messageEmbeds := []*discordgo.MessageEmbedField{
		{Name: "WarnAt", Value: warnAt},
		{Name: "Description", Value: description},
	}

	additionalMessageFields := []*discordgo.MessageEmbed{}

	for i := 0; i < len(fields); i += 2 {
		additionalMessageFields = append(additionalMessageFields, &discordgo.MessageEmbed{
			Title:       fields[i],
			Description: fields[i+1],
		})
	}

	err := n.SendOnDiscord(warnString, warnColorINT, n.config.Warn, messageEmbeds, additionalMessageFields)
	if err != nil {
		return fmt.Errorf("❌ Failed to report warn on Discord: %s", err)
	}

	return nil
}

func (n *Notifier) SendOnDiscord(text string, messageColor int, channelConfig *DiscordChannelConfig, messageEmbeds []*discordgo.MessageEmbedField, additionalMessageFields []*discordgo.MessageEmbed) error {
	mentions := generateMentions(channelConfig.Mentions)

	// Create the Discord message with embeds that we will send to the channel
	message := &discordgo.MessageSend{
		Content: text,
		Embed: &discordgo.MessageEmbed{
			Title:       mentions,
			Description: text,
			Color:       messageColor,
			Fields:      messageEmbeds,
			Timestamp:   time.Now().Format(time.RFC3339),
		},
	}

	// Send the message as a new thread
	threadMessage, err := n.session.ChannelMessageSendComplex(channelConfig.ChannelID, message)
	if err != nil {
		return fmt.Errorf("failed to send message : %s", err)
	}
	// Set the reply message as a reply to the original message
	replyReference := &discordgo.MessageReference{
		MessageID: threadMessage.ID,
		ChannelID: threadMessage.ChannelID,
		GuildID:   threadMessage.GuildID,
	}
	// Send additional message fields as replies in the thread
	for _, embed := range additionalMessageFields {
		embed.Timestamp = time.Now().Format(time.RFC3339)
		reply := &discordgo.MessageSend{
			Embed:     embed,
			Reference: replyReference,
		}
		_, err := n.session.ChannelMessageSendComplex(channelConfig.ChannelID, reply)
		if err != nil {
			continue
		}
	}

	return nil
}

func isConfigured(channelConfig *DiscordChannelConfig) bool {
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
		str += fmt.Sprintf("@%s ", user)
	}

	return str
}

package notifier

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	errorColor      = "#FF0000"
	successColor    = "#36a64f"
	warnColor       = "#fceea7"
	warnColorINT    = 16580711
	successColorINT = 3559039
	errorColorINT   = 16711680
)

var standardNotifier *Notifier

type Notifier struct {
	session *discordgo.Session
	config  *Config
}

func Init(config *Config) *Notifier {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		logrus.Fatal("discord notifier can not be initialized error: %s", err)
	}
	standardNotifier = &Notifier{
		session: session,
		config:  config,
	}

	return standardNotifier
}

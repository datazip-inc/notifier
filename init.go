package notifier

import (
	"github.com/slack-go/slack"
)

const (
	errorColor   = "#FF0000"
	successColor = "#36a64f"
	warnColor    = "#fceea7"
)

var standardNotifier *Notifier

type Notifier struct {
	slackClient *slack.Client
	config      *Config
}

func Init(config *Config) *Notifier {
	standardNotifier = &Notifier{
		slackClient: slack.New(config.Token, slack.OptionDebug(false)),
		config:      config,
	}

	return standardNotifier
}

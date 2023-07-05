package notifier

type Config struct {
	Token   string              `yaml:"token"`
	Default *SlackChannelConfig `yaml:"default"`
	Debug   *SlackChannelConfig `yaml:"debug"`
}

type SlackChannelConfig struct {
	ChannelID string   `yaml:"id"`
	Mentions  []string `yaml:"mentions"` // mentions are expected to have @ excluded before usernames
}

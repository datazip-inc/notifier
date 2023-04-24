package notifier

type Config struct {
	Token   string                `yaml:"token"`
	Success *DiscordChannelConfig `yaml:"success"`
	Error   *DiscordChannelConfig `yaml:"error"`
	Warn    *DiscordChannelConfig `yaml:"warn"`
}

type DiscordChannelConfig struct {
	ChannelID string   `yaml:"id"`
	Mentions  []string `yaml:"mentions"` // mentions are expected to have @ excluded before usernames
}

// type Config struct {
// 	Token   string              `yaml:"token"`
// 	Success *SlackChannelConfig `yaml:"success"`
// 	Error   *SlackChannelConfig `yaml:"error"`
// 	Warn    *SlackChannelConfig `yaml:"warn"`
// }

// slack channel config will be deprecated
// type SlackChannelConfig struct {
// 	ChannelID string   `yaml:"id"`
// 	Mentions  []string `yaml:"mentions"` // mentions are expected to have @ excluded before usernames
// }

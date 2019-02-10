package cli

import "github.com/jessevdk/go-flags"

type CLI struct {
	TelegramBotToken string `short:"t" long:"token" description:"token for authorization to telegram bot"`
	Debug            bool   `long:"debug" description:"enable telegram debug"`
	Backend          string `short:"b" long:"backend" description:"toussaint backend url" default:"127.0.0.1:8080"`
}

func ParseCLI() (*CLI, error) {
	parser := flags.NewNamedParser("toussaint-telegram", flags.Default)
	var params CLI
	_, err := parser.AddGroup("Common", "", &params)
	if err != nil {
		return nil, err
	}

	_, err = parser.Parse()
	if err != nil {
		return nil, err
	}

	return &params, nil
}

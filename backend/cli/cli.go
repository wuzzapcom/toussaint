package cli

import "github.com/jessevdk/go-flags"

type CLIParams struct {
	TimeBetweenUpdatesStr string `long:"between" description:"time between game updates" default:"24h"`
	UpdateTimeStr         string `long:"time" description:"time of first update in format hh:mm" default:"02:00"`
	Host                  string `long:"host" description:"host that srv listen to" default:"127.0.0.1"`
	Post                  int    `long:"port" description:"port that srv listen to" default:"8080"`
	Debug                 bool   `short:"d" long:"debug" description:"run backend in debug mode with additional test endpoints"`
}

func ParseCLI() (*CLIParams, error) {
	parser := flags.NewNamedParser("toussaint", flags.Default)
	var params CLIParams
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

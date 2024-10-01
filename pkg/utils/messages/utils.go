package messages

import (
	"fmt"
	"github.com/pterm/pterm"
)

type SuccessMsgParams struct {
	Domain       string
	HTTPSEnabled bool
	DataDir      string
	IndexFile    string
	StaticDir    string
	ConfigFile   string
	EnvFile      string
	ServiceFile  string
	Service      string
	RelayName    string
	GitHubLink   string
}

func SuccessMessages(successMsgParams *SuccessMsgParams) {
	pterm.Println()
	pterm.Println(pterm.Green("The installation is complete."))

	pterm.Println()
	pterm.Println(pterm.Cyan("You can access your relay at:"))
	if successMsgParams.HTTPSEnabled {
		pterm.Println(pterm.Magenta("wss://" + successMsgParams.Domain))
	} else {
		pterm.Println(pterm.Magenta("ws://" + successMsgParams.Domain))
	}

	pterm.Println()
	pterm.Println(pterm.Cyan("Your relay's data directory is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.DataDir))

	if successMsgParams.IndexFile != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's index.html file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.IndexFile))
	}

	if successMsgParams.StaticDir != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's static directory is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.StaticDir))
	}

	if successMsgParams.ConfigFile != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's config file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.ConfigFile))
	}

	if successMsgParams.EnvFile != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's environment file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.EnvFile))
	}

	pterm.Println()
	pterm.Println(pterm.Cyan("Your relay's service file is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.ServiceFile))

	pterm.Println()
	pterm.Println(pterm.Cyan("To check the status of your relay run:"))
	pterm.Println(pterm.Magenta("systemctl status " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Cyan("To reload the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl reload " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Cyan("To restart the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl restart " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Cyan(fmt.Sprintf("%s GitHub", successMsgParams.RelayName)))
	pterm.Println(pterm.Magenta(successMsgParams.GitHubLink))
}

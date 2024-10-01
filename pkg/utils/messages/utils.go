package messages

import (
	"fmt"
	"github.com/pterm/pterm"
)

type SuccessMsgParams struct {
	Domain      string
	DataDir     string
	IndexFile   string
	StaticDir   string
	ConfigFile  string
	EnvFile     string
	ServiceFile string
	Service     string
	RelayName   string
	GitHubLink  string
}

func SuccessMessages(successMsgParams *SuccessMsgParams) {
	pterm.Println()
	pterm.Println(pterm.Magenta("The installation is complete."))

	pterm.Println()
	pterm.Println(pterm.Magenta("You can access your relay at:"))
	pterm.Println(pterm.Magenta("wss://" + successMsgParams.Domain))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's data directory is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.DataDir))

	if successMsgParams.IndexFile != "" {
		pterm.Println()
		pterm.Println(pterm.Magenta("Your relay's index.html file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.IndexFile))
	}

	if successMsgParams.StaticDir != "" {
		pterm.Println()
		pterm.Println(pterm.Magenta("Your relay's static directory is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.StaticDir))
	}

	if successMsgParams.ConfigFile != "" {
		pterm.Println()
		pterm.Println(pterm.Magenta("Your relay's config file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.ConfigFile))
	}

	if successMsgParams.EnvFile != "" {
		pterm.Println()
		pterm.Println(pterm.Magenta("Your relay's environment file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.EnvFile))
	}

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's service file is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.ServiceFile))

	pterm.Println()
	pterm.Println(pterm.Magenta("To check the status of your relay run:"))
	pterm.Println(pterm.Magenta("systemctl status " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Magenta("To reload the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl reload " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Magenta("To restart the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl restart " + successMsgParams.Service))

	pterm.Println()
	pterm.Println(pterm.Magenta(fmt.Sprintf("%s GitHub", successMsgParams.RelayName)))
	pterm.Println(pterm.Magenta(successMsgParams.GitHubLink))
}

package messages

import (
	"fmt"
	"github.com/pterm/pterm"
)

type SuccessMsgParams struct {
	Domain               string
	HTTPSEnabled         bool
	DataDirPath          string
	IndexFilePath        string
	StaticDirPath        string
	ConfigFilePath       string
	PluginFilePath       string
	NginxConfigFilePath  string
	BinaryFilePath       string
	BinaryPluginFilePath string
	EnvFilePath          string
	ServiceFilePath      string
	ServiceName          string
	RelayName            string
	GitHubLink           string
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
	pterm.Println(pterm.Magenta(successMsgParams.DataDirPath))

	if successMsgParams.IndexFilePath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's index.html file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.IndexFilePath))
	}

	if successMsgParams.StaticDirPath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's static directory is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.StaticDirPath))
	}

	if successMsgParams.ConfigFilePath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's config file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.ConfigFilePath))
	}

	if successMsgParams.PluginFilePath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's plugin file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.PluginFilePath))
	}

	pterm.Println()
	pterm.Println(pterm.Cyan("Your relay's nginx config file is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.NginxConfigFilePath))

	pterm.Println()
	pterm.Println(pterm.Cyan("Your relay's binary file is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.BinaryFilePath))

	if successMsgParams.BinaryPluginFilePath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's binary plugin file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.BinaryPluginFilePath))
	}

	if successMsgParams.EnvFilePath != "" {
		pterm.Println()
		pterm.Println(pterm.Cyan("Your relay's environment file is located here:"))
		pterm.Println(pterm.Magenta(successMsgParams.EnvFilePath))
	}

	pterm.Println()
	pterm.Println(pterm.Cyan("Your relay's service file is located here:"))
	pterm.Println(pterm.Magenta(successMsgParams.ServiceFilePath))

	pterm.Println()
	pterm.Println(pterm.Cyan("To check the status of your relay run:"))
	pterm.Println(pterm.Magenta("systemctl status " + successMsgParams.ServiceName))

	pterm.Println()
	pterm.Println(pterm.Cyan("To reload the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl reload " + successMsgParams.ServiceName))

	pterm.Println()
	pterm.Println(pterm.Cyan("To restart the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl restart " + successMsgParams.ServiceName))

	pterm.Println()
	pterm.Println(pterm.Cyan(fmt.Sprintf("%s GitHub", successMsgParams.RelayName)))
	pterm.Println(pterm.Magenta(successMsgParams.GitHubLink))
}

package wot_relay

import (
	"github.com/pterm/pterm"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/wot-relay"
	const indexFile = "/etc/wot-relay/templates/index.html"
	const staticDir = "/etc/wot-relay/templates/static"
	const envFile = "/etc/systemd/system/wot-relay.env"
	const serviceFile = "/etc/systemd/system/wot-relay.service"
	const service = "wot-relay"
	const githubLink = "https://github.com/bitvora/wot-relay"

	pterm.Println()
	pterm.Println(pterm.Magenta("The installation is complete."))

	pterm.Println()
	pterm.Println(pterm.Magenta("You can access your relay at:"))
	pterm.Println(pterm.Magenta("wss://" + domain))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's data directory is located here:"))
	pterm.Println(pterm.Magenta(dataDir))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's index.html file is located here:"))
	pterm.Println(pterm.Magenta(indexFile))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's static directory is located here:"))
	pterm.Println(pterm.Magenta(staticDir))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's environment file is located here:"))
	pterm.Println(pterm.Magenta(envFile))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's service file is located here:"))
	pterm.Println(pterm.Magenta(serviceFile))

	pterm.Println()
	pterm.Println(pterm.Magenta("To check the status of your relay run:"))
	pterm.Println(pterm.Magenta("systemctl status " + service))

	pterm.Println()
	pterm.Println(pterm.Magenta("To reload the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl reload " + service))

	pterm.Println()
	pterm.Println(pterm.Magenta("To restart the relay service run:"))
	pterm.Println(pterm.Magenta("systemctl restart " + service))

	pterm.Println()
	pterm.Println(pterm.Magenta("WoT Relay GitHub"))
	pterm.Println(pterm.Magenta(githubLink))
}

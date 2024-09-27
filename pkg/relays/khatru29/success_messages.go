package khatru29

import (
	"github.com/pterm/pterm"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/khatru29"
	const envFile = "/etc/systemd/system/khatru29.env"
	const serviceFile = "/etc/systemd/system/khatru29.service"
	const service = "khatru29"
	const githubLink = "https://github.com/fiatjaf/relay29/tree/master"

	pterm.Println()
	pterm.Println(pterm.Magenta("The installation is complete."))

	pterm.Println()
	pterm.Println(pterm.Magenta("You can access your relay at:"))
	pterm.Println(pterm.Magenta("wss://" + domain))

	pterm.Println()
	pterm.Println(pterm.Magenta("Your relay's data directory is located here:"))
	pterm.Println(pterm.Magenta(dataDir))

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
	pterm.Println(pterm.Magenta("Khatru29 GitHub"))
	pterm.Println(pterm.Magenta(githubLink))
}

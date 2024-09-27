package khatru_pyramid

import (
	"github.com/pterm/pterm"
)

func SuccessMessages(domain string) {
	const dataDir = "/var/lib/khatru-pyramid"
	const envFile = "/etc/systemd/system/khatru-pyramid.env"
	const serviceFile = "/etc/systemd/system/khatru-pyramid.service"
	const service = "khatru-pyramid"
	const githubLink = "https://github.com/github-tijlxyz/khatru-pyramid"

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
	pterm.Println(pterm.Magenta("Khatru Pyramid GitHub"))
	pterm.Println(pterm.Magenta(githubLink))
}

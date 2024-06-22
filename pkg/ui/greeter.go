package ui

import (
	"github.com/pterm/pterm"
)

func Greet() {
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		pterm.Magenta("\nWelcome to Relay Wizard 🪄") + pterm.Gray("\nInstall and manage your relays with ease!") + pterm.Gray("\nv0.0.1"))
}

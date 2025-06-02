package cmd

import (
	"fmt"
	"github.com/nodetec/rwz/pkg/manager"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/relays/khatru29"
	"github.com/nodetec/rwz/pkg/relays/khatru_pyramid"
	"github.com/nodetec/rwz/pkg/relays/nostr_rs_relay"
	"github.com/nodetec/rwz/pkg/relays/strfry"
	"github.com/nodetec/rwz/pkg/relays/strfry29"
	"github.com/nodetec/rwz/pkg/relays/wot_relay"
	"github.com/nodetec/rwz/pkg/ui"
	"github.com/nodetec/rwz/pkg/utils/users"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure your Nostr relay",
	Long:  `Install and configure your Nostr relay, including package installation, firewall setup, Nginx configuration, SSL/TLS certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {
		ThemeDefault := pterm.ThemeDefault

		ui.Greet()

		relayDomain, _ := pterm.DefaultInteractiveTextInput.Show("Relay domain name")
		pterm.Println()

		// Supported relay options
		options := []string{khatru_pyramid.RelayName, nostr_rs_relay.RelayName, strfry.RelayName, wot_relay.RelayName, khatru29.RelayName, strfry29.RelayName}

		// Use PTerm's interactive select feature to present the options to the user and capture their selection
		// The Show() method displays the options and waits for the user's input
		relaySelector := pterm.InteractiveSelectPrinter{
			TextStyle:     &ThemeDefault.PrimaryStyle,
			DefaultText:   "Please select an option",
			Options:       []string{},
			OptionStyle:   &ThemeDefault.DefaultText,
			DefaultOption: "",
			MaxHeight:     6,
			Selector:      ">",
			SelectorStyle: &ThemeDefault.SecondaryStyle,
			Filter:        true,
		}

		selectedRelayOption, _ := relaySelector.WithOptions(options).Show()

		// Display the selected option to the user with a green color for emphasis
		pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedRelayOption))

		var privKey string
		var pubKey string
		if selectedRelayOption == khatru_pyramid.RelayName || selectedRelayOption == nostr_rs_relay.RelayName || selectedRelayOption == strfry.RelayName || selectedRelayOption == wot_relay.RelayName || selectedRelayOption == strfry29.RelayName {
			pterm.Println()
			pubKey, _ = pterm.DefaultInteractiveTextInput.Show("Public key (hex not npub)")
		}

		if selectedRelayOption == khatru29.RelayName || selectedRelayOption == strfry29.RelayName {
			pterm.Println()
			privKeyInput := pterm.DefaultInteractiveTextInput.WithMask("*")
			privKey, _ = privKeyInput.Show("Private key (hex not nsec)")
		}

		var relayContact string
		if selectedRelayOption == wot_relay.RelayName {
			pterm.Println()
			pterm.Println(pterm.Yellow("If you leave the relay contact information empty, then the relay's public key will be used."))

			pterm.Println()
			relayContact, _ = pterm.DefaultInteractiveTextInput.Show("Email address/Public key (hex not npub)")
		} else {
			pterm.Println()
			pterm.Println(pterm.Yellow("Leave email empty if you don't want to provide relay contact information."))

			pterm.Println()
			relayContact, _ = pterm.DefaultInteractiveTextInput.Show("Email address")
			if relayContact != "" {
				relayContact = fmt.Sprintf("mailto:%s", relayContact)
			}
		}

		pterm.Println()
		pterm.Println(pterm.Yellow("If you make a mistake, you can always re-run this installer."))
		pterm.Println()

		// Install necessary packages using APT
		manager.AptInstallPackages(selectedRelayOption)

		// Configure the firewall
		network.ConfigureFirewall()

		// Configure the intrusion detection system
		network.ConfigureIntrusionDetection()

		// Configure Nginx
		network.ConfigureNginx()

		// Create relay user
		var relayUser string
		pterm.Println()
		pterm.Println(pterm.Cyan("Create a user for the relay."))

		pterm.Println()
		userInput := pterm.DefaultInteractiveTextInput.WithDefaultValue(relays.DefaultUser)
		relayUser, _ = userInput.Show("Relay user")

		pterm.Println()
		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking if '%s' user exists...", relayUser))
		if !users.UserExists(relayUser) {
			spinner.UpdateText(fmt.Sprintf("Creating '%s' user...", relayUser))
			users.CreateUser(relayUser, true)
			spinner.Success(fmt.Sprintf("Created '%s' user", relayUser))
		} else {
			spinner.Success(fmt.Sprintf("'%s' user already exists", relayUser))
		}

		if selectedRelayOption == khatru_pyramid.RelayName {
			khatru_pyramid.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == nostr_rs_relay.RelayName {
			nostr_rs_relay.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == strfry.RelayName {
			strfry.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == wot_relay.RelayName {
			wot_relay.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == khatru29.RelayName {
			khatru29.Install(relayDomain, privKey, relayContact, relayUser)
		} else if selectedRelayOption == strfry29.RelayName {
			strfry29.Install(relayDomain, pubKey, privKey, relayContact, relayUser)
		}

		pterm.Println()
		pterm.Println(pterm.Cyan("Join the NODE-TEC Discord to get support:"))
		pterm.Println(pterm.Magenta("https://discord.gg/J9gRK5pbWb"))
		pterm.Println()
		pterm.Println(pterm.Cyan("We plan to use relay groups for support in the future..."))

		pterm.Println()
		pterm.Println(pterm.Magenta("You can re-run this installer with `rwz install`."))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

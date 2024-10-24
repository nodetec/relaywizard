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
		if selectedRelayOption == khatru_pyramid.RelayName || selectedRelayOption == nostr_rs_relay.RelayName || selectedRelayOption == wot_relay.RelayName {
			pterm.Println()
			pubKey, _ = pterm.DefaultInteractiveTextInput.Show("Public key (hex not npub)")
		} else if selectedRelayOption == khatru29.RelayName || selectedRelayOption == strfry29.RelayName {
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

		// Step 1: Install necessary packages using APT
		manager.AptInstallPackages(selectedRelayOption)

		// Step 2: Configure the firewall
		network.ConfigureFirewall()

		// Setp 3: Create relay user
		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Checking if '%s' user exists...", relays.User))
		if !users.UserExists(relays.User) {
			spinner.UpdateText(fmt.Sprintf("Creating '%s' user...", relays.User))
			users.CreateUser(relays.User, true)
			spinner.Success(fmt.Sprintf("Created '%s' user.", relays.User))
		} else {
			spinner.Success(fmt.Sprintf("'%s' user already exists.", relays.User))
		}

		if selectedRelayOption == khatru_pyramid.RelayName {
			// Step 4: Configure Nginx for HTTP
			khatru_pyramid.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				khatru_pyramid.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			khatru_pyramid.InstallRelayBinary()

			// Step 8: Set up the relay service
			khatru_pyramid.SetupRelayService(relayDomain, pubKey, relayContact)

			// Step 9: Show success messages
			khatru_pyramid.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == nostr_rs_relay.RelayName {
			// Step 4: Configure Nginx for HTTP
			nostr_rs_relay.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				nostr_rs_relay.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			nostr_rs_relay.InstallRelayBinary()

			// Step 8: Set up the relay service
			nostr_rs_relay.SetupRelayService(relayDomain, pubKey, relayContact, httpsEnabled)

			// Step 9: Show success messages
			nostr_rs_relay.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == strfry.RelayName {
			// Step 4: Configure Nginx for HTTP
			strfry.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				strfry.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			strfry.InstallRelayBinary()

			// Step 8: Set up the relay service
			strfry.SetupRelayService(relayDomain, relayContact)

			// Step 9: Show success messages
			strfry.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == wot_relay.RelayName {
			// Step 4: Configure Nginx for HTTP
			wot_relay.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				wot_relay.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			wot_relay.InstallRelayBinary()

			// Step 8: Set up the relay service
			wot_relay.SetupRelayService(relayDomain, pubKey, relayContact, httpsEnabled)

			// Step 9: Show success messages
			wot_relay.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == khatru29.RelayName {
			// Step 4: Configure Nginx for HTTP
			khatru29.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				khatru29.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			khatru29.InstallRelayBinary()

			// Step 8: Set up the relay service
			khatru29.SetupRelayService(relayDomain, privKey, relayContact)

			// Step 9: Show success messages
			khatru29.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == strfry29.RelayName {
			// Step 4: Configure Nginx for HTTP
			strfry29.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL/TLS certificates
			httpsEnabled := network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 6: Configure Nginx for HTTPS
				strfry29.ConfigureNginxHttps(relayDomain)
			}

			// Step 7: Download and install the relay binary
			strfry29.InstallRelayBinary()

			// Step 8: Set up the relay service
			strfry29.SetupRelayService(relayDomain, privKey, relayContact)

			// Step 9: Show success messages
			strfry29.SuccessMessages(relayDomain, httpsEnabled)
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

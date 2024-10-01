package cmd

import (
	"github.com/nodetec/rwz/pkg/manager"
	"github.com/nodetec/rwz/pkg/network"
	"github.com/nodetec/rwz/pkg/relays/khatru29"
	"github.com/nodetec/rwz/pkg/relays/khatru_pyramid"
	"github.com/nodetec/rwz/pkg/relays/strfry"
	"github.com/nodetec/rwz/pkg/relays/wot_relay"
	"github.com/nodetec/rwz/pkg/ui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure your Nostr relay",
	Long:  `Install and configure your Nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {

		ui.Greet()

		relayDomain, _ := pterm.DefaultInteractiveTextInput.Show("Relay domain name")
		pterm.Println()

		// Supported relay options
		options := []string{"Khatru Pyramid", "strfry", "Khatru29", "WoT Relay"}

		// Use PTerm's interactive select feature to present the options to the user and capture their selection
		// The Show() method displays the options and waits for the user's input
		selectedRelayOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		// Display the selected option to the user with a green color for emphasis
		pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedRelayOption))

		var privkey string
		var pubkey string
		if selectedRelayOption == "Khatru Pyramid" || selectedRelayOption == "WoT Relay" {
			pterm.Println()
			pubkey, _ = pterm.DefaultInteractiveTextInput.Show("Public key (hex not npub)")
		} else if selectedRelayOption == "Khatru29" {
			pterm.Println()
			privkey, _ = pterm.DefaultInteractiveTextInput.Show("Private key (hex not nsec)")
		}

		pterm.Println()
		pterm.Println(pterm.Yellow("If you make a mistake, you can always re-run this installer."))
		pterm.Println()

		// Step 1: Install necessary packages using APT
		manager.AptInstallPackages()

		if selectedRelayOption == "Khatru Pyramid" {
			// Step 2: Configure the firewall
			network.ConfigureFirewall()

			// Step 3: Configure Nginx for HTTP
			khatru_pyramid.ConfigureNginxHttp(relayDomain)

			// Step 4: Get SSL certificates
			var httpsEnabled = network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 5: Configure Nginx for HTTPS
				khatru_pyramid.ConfigureNginxHttps(relayDomain)
			}

			// Step 6: Download and install the relay binary
			khatru_pyramid.InstallRelayBinary()

			// Step 7: Set up the relay service
			khatru_pyramid.SetupRelayService(relayDomain, pubkey)

			// Step 8: Show success messages
			khatru_pyramid.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == "strfry" {
			// Step 2: Configure the firewall
			network.ConfigureFirewall()

			// Step 3: Configure Nginx for HTTP
			strfry.ConfigureNginxHttp(relayDomain)

			// Step 4: Get SSL certificates
			var httpsEnabled = network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 5: Configure Nginx for HTTPS
				strfry.ConfigureNginxHttps(relayDomain)
			}

			// Step 6: Download and install the relay binary
			strfry.InstallRelayBinary()

			// Step 7: Set up the relay service
			strfry.SetupRelayService(relayDomain)

			// Step 8: Show success messages
			strfry.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == "Khatru29" {
			// Step 2: Configure the firewall
			network.ConfigureFirewall()

			// Step 3: Configure Nginx for HTTP
			khatru29.ConfigureNginxHttp(relayDomain)

			// Step 4: Get SSL certificates
			var httpsEnabled = network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 5: Configure Nginx for HTTPS
				khatru29.ConfigureNginxHttps(relayDomain)
			}

			// Step 6: Download and install the relay binary
			khatru29.InstallRelayBinary()

			// Step 7: Set up the relay service
			khatru29.SetupRelayService(relayDomain, privkey)

			// Step 8: Show success messages
			khatru29.SuccessMessages(relayDomain, httpsEnabled)
		} else if selectedRelayOption == "WoT Relay" {
			// Step 2: Configure the firewall
			network.ConfigureFirewall()

			// Step 3: Configure Nginx for HTTP
			wot_relay.ConfigureNginxHttp(relayDomain)

			// Step 4: Get SSL certificates
			var httpsEnabled = network.GetCertificates(relayDomain)
			if httpsEnabled {
				// Step 5: Configure Nginx for HTTPS
				wot_relay.ConfigureNginxHttps(relayDomain)
			}

			// Step 6: Download and install the relay binary
			wot_relay.InstallRelayBinary()

			// Step 7: Set up the relay service
			wot_relay.SetupRelayService(relayDomain, pubkey, httpsEnabled)

			// Step 8: Show success messages
			wot_relay.SuccessMessages(relayDomain, httpsEnabled)
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

package cmd

import (
	"github.com/nodetec/relaywiz/pkg/manager"
	"github.com/nodetec/relaywiz/pkg/network"
	"github.com/nodetec/relaywiz/pkg/relays/khatru29"
	"github.com/nodetec/relaywiz/pkg/relays/khatru_pyramid"
	"github.com/nodetec/relaywiz/pkg/relays/strfry"
	"github.com/nodetec/relaywiz/pkg/ui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure the nostr relay",
	Long:  `Install and configure the nostr relay, including package installation, nginx configuration, firewall setup, SSL certificates, and starting the relay service.`,
	Run: func(cmd *cobra.Command, args []string) {

		ui.Greet()

		relayDomain, _ := pterm.DefaultInteractiveTextInput.Show("Relay domain name")
		pterm.Println()
		pterm.Println(pterm.Yellow("Leave email empty if you don't want to receive notifications from Let's Encrypt about your SSL cert."))
		pterm.Println()
		ssl_email, _ := pterm.DefaultInteractiveTextInput.Show("Email address")
		pterm.Println()

		// Supported relay options
		options := []string{"Khatru Pyramid", "strfry", "khatru29"}

		// Use PTerm's interactive select feature to present the options to the user and capture their selection
		// The Show() method displays the options and waits for the user's input
		selectedRelayOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

		// Display the selected option to the user with a green color for emphasis
		pterm.Info.Printfln("Selected option: %s", pterm.Green(selectedRelayOption))

		var privkey string
		var pubkey string
		if selectedRelayOption == "Khatru Pyramid" {
			pterm.Println()
			pubkey, _ = pterm.DefaultInteractiveTextInput.Show("Public key (hex not npub)")
		} else if selectedRelayOption == "khatru29" {
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
			var shouldContinue = network.GetCertificates(relayDomain, ssl_email)
			if !shouldContinue {
				return
			}

			// Step 5: Configure Nginx for HTTPS
			khatru_pyramid.ConfigureNginxHttps(relayDomain)

			// Step 6: Download and install the relay binary
			khatru_pyramid.InstallRelayBinary()

			// Step 7: Set up the relay service
			khatru_pyramid.SetupRelayService(relayDomain, pubkey)

			// Step 8: Show success messages
			khatru_pyramid.SuccessMessages(relayDomain)
		} else if selectedRelayOption == "strfry" {
			// Step 2: Install necessary strfry package dependencies
			strfry.AptInstallDependencies()

			// Step 3: Configure the firewall
			network.ConfigureFirewall()

			// Step 4: Configure Nginx for HTTP
			strfry.ConfigureNginxHttp(relayDomain)

			// Step 5: Get SSL certificates
			var shouldContinue = network.GetCertificates(relayDomain, ssl_email)
			if !shouldContinue {
				return
			}

			// Step 6: Configure Nginx for HTTPS
			strfry.ConfigureNginxHttps(relayDomain)

			// Step 7: Download and install the relay binary
			strfry.InstallRelayBinary()

			// Step 8: Set up the relay service
			strfry.SetupRelayService(relayDomain)

			// Step 9: Show success messages
			strfry.SuccessMessages(relayDomain)
		} else if selectedRelayOption == "khatru29" {
			// Step 2: Configure the firewall
			network.ConfigureFirewall()

			// Step 3: Configure Nginx for HTTP
			khatru29.ConfigureNginxHttp(relayDomain)

			// Step 4: Get SSL certificates
			var shouldContinue = network.GetCertificates(relayDomain, ssl_email)
			if !shouldContinue {
				return
			}

			// Step 5: Configure Nginx for HTTPS
			khatru29.ConfigureNginxHttps(relayDomain)

			// Step 6: Download and install the relay binary
			khatru29.InstallRelayBinary()

			// Step 7: Set up the relay service
			khatru29.SetupRelayService(relayDomain, privkey)

			// Step 8: Show success messages
			khatru29.SuccessMessages(relayDomain)
		}

		pterm.Println()
		pterm.Println(pterm.Magenta("Join the NODE-TEC Discord to get support:"))
		pterm.Println(pterm.Magenta("https://discord.gg/J9gRK5pbWb"))
		pterm.Println()
		pterm.Println(pterm.Magenta("We plan to use relay groups for support in the future..."))

		pterm.Println()
		pterm.Println(pterm.Magenta("You can re-run this installer with `relaywiz install`."))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

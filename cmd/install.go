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
	"github.com/nodetec/rwz/pkg/utils/files"
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

		// TODO
		// Add check here after getting the domain

		// Supported relay options
		options := []string{relays.KhatruPyramidRelayName, relays.NostrRsRelayName, relays.StrfryRelayName, relays.WotRelayName, relays.Khatru29RelayName, relays.Strfry29RelayName}

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
		if selectedRelayOption == relays.KhatruPyramidRelayName || selectedRelayOption == relays.NostrRsRelayName || selectedRelayOption == relays.StrfryRelayName || selectedRelayOption == relays.WotRelayName || selectedRelayOption == relays.Strfry29RelayName {
			pterm.Println()
			pubKey, _ = pterm.DefaultInteractiveTextInput.Show("Public key (hex not npub)")
		}

		if selectedRelayOption == relays.Khatru29RelayName || selectedRelayOption == relays.Strfry29RelayName {
			pterm.Println()
			privKeyInput := pterm.DefaultInteractiveTextInput.WithMask("*")
			privKey, _ = privKeyInput.Show("Private key (hex not nsec)")
		}

		var relayContact string
		if selectedRelayOption == relays.WotRelayName {
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

		pterm.Println()
		pterm.Println(pterm.Yellow("Warning: Relay Wizard SSH defaults will not be applied if the current sshd configuration overrides them."))
		pterm.Printfln(pterm.Yellow("If issues occur, try checking the following locations %s and %s"), network.SSHDConfigFilePath, network.SSHDConfigDDirPath)

		prompt := pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Configure remote access through SSH using Relay Wizard defaults?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"yes", "no"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		pterm.Println()
		result, _ := prompt.Show()
		pterm.Println()

		if result == "yes" {
			network.ConfigureRemoteAccess()
		} else {
			files.RemoveFile(network.SSHDConfigDRWZConfigFilePath)
		}

		// Configure the firewall
		network.ConfigureFirewall()

		// Configure the intrusion detection system
		network.ConfigureIntrusionDetection()

		// Configure Nginx
		network.ConfigureNginx()

		// Create relay user
		var relayUser string
		pterm.Println()
		pterm.Println(pterm.LightCyan("Create a user for the relay."))

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

		if selectedRelayOption == relays.KhatruPyramidRelayName {
			khatru_pyramid.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == relays.NostrRsRelayName {
			nostr_rs_relay.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == relays.StrfryRelayName {
			strfry.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == relays.WotRelayName {
			wot_relay.Install(relayDomain, pubKey, relayContact, relayUser)
		} else if selectedRelayOption == relays.Khatru29RelayName {
			khatru29.Install(relayDomain, privKey, relayContact, relayUser)
		} else if selectedRelayOption == relays.Strfry29RelayName {
			strfry29.Install(relayDomain, pubKey, privKey, relayContact, relayUser)
		}

		pterm.Println()
		pterm.Println(pterm.LightCyan("Join the NODE-TEC Discord to get support:"))
		pterm.Println(pterm.Magenta("https://discord.gg/J9gRK5pbWb"))
		pterm.Println()
		pterm.Println(pterm.LightCyan("We plan to use relay groups for support in the future..."))

		pterm.Println()
		pterm.Println(pterm.Magenta("You can re-run this installer with `rwz install`."))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

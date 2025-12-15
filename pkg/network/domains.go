package network

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/directories"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to check if the provided relay domain has already been used
func CheckDomainUsage(currentUsername, relayDomain, nginxConfigFilePath string) {
	spinner, _ := pterm.DefaultSpinner.Start("Checking relay domain availability...")

	var nginxTOutput []byte
	var err error
	if currentUsername == relays.RootUser {
		files.RemoveFile(nginxConfigFilePath)
		nginxTOutput, err = exec.Command("nginx", "-T").CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to test and display Nginx configuration: %v", err)
			os.Exit(1)
		}
	} else {
		files.RemoveFileUsingLinux(currentUsername, nginxConfigFilePath)
		nginxTOutput, err = exec.Command("sudo", "nginx", "-T").CombinedOutput()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to test and display Nginx configuration: %v", err)
			os.Exit(1)
		}
	}

	nginxTestAndConfigurationOutput := string(nginxTOutput)

	nginxConfigurationTestSuccessOutput := "nginx: configuration file /etc/nginx/nginx.conf test is successful"

	_, completeNginxConfiguration, found := strings.Cut(nginxTestAndConfigurationOutput, nginxConfigurationTestSuccessOutput)

	if !found {
		pterm.Println()
		pterm.Error.Printfln("Failed to find Nginx configuration success output")
		os.Exit(1)
	}

	relayDomainNameFound := strings.Contains(completeNginxConfiguration, fmt.Sprintf("server_name %s", relayDomain))

	if relayDomainNameFound {
		pterm.Println()
		pterm.Error.Printfln(fmt.Sprintf("The domain %s is already being used in another Nginx configuration file.", relayDomain))
		pterm.Error.Printfln("Try using a different domain for the relay or updating/removing the conflicting Nginx configuration file.")
		pterm.Error.Printfln(fmt.Sprintf("Nginx configuration files can be found here: %s", NginxConfDirPath))
		os.Exit(1)
	}

	spinner.Success("Relay domain available")

	pterm.Println()
	spinner, _ = pterm.DefaultSpinner.Start("Checking for existing relay site...")

	relayDomainDirPath := fmt.Sprintf("%s/%s", WWWDirPath, relayDomain)

	if directories.DirExists(relayDomainDirPath) {
		spinner.Info("Relay site found...")

		pterm.Println()
		pterm.Printfln(pterm.LightCyan("If you haven't modified any content or installed a different type of relay in %s, then you can safely overwrite the relay site."), relayDomainDirPath)
		pterm.Printfln(pterm.LightCyan("If you have installed a different type of relay in %s and you're no longer using that relay, then you can also safely overwrite the relay site."), relayDomainDirPath)
		pterm.Println()

		ThemeDefault := pterm.ThemeDefault

		prompt := pterm.InteractiveContinuePrinter{
			DefaultValueIndex: 0,
			DefaultText:       "Overwrite relay site?",
			TextStyle:         &ThemeDefault.PrimaryStyle,
			Options:           []string{"no", "yes"},
			OptionsStyle:      &ThemeDefault.SuccessMessageStyle,
			SuffixStyle:       &ThemeDefault.SecondaryStyle,
			Delimiter:         ": ",
		}

		result, _ := prompt.Show()

		if result == "no" {
			pterm.Println()
			pterm.Info.Printfln("Canceling installation...")
			os.Exit(1)
		} else if result == "yes" {
			pterm.Println()
			pterm.Println(pterm.Yellow("Warning: Are you sure you want to overwrite your existing relay site?"))
			pterm.Printfln(pterm.Yellow("If you select 'yes', then the following relay site will be overwritten: %s"), relayDomainDirPath)
			pterm.Println()

			result, _ := prompt.Show()

			if result == "no" {
				pterm.Println()
				pterm.Info.Printfln("Canceling installation...")
				os.Exit(1)
			} else if result == "yes" {
				if currentUsername == relays.RootUser {
					directories.RemoveDirectory(relayDomainDirPath)
				} else {
					directories.RemoveDirectoryUsingLinux(currentUsername, relayDomainDirPath)
				}
				pterm.Println()
				pterm.Println(pterm.LightCyan("Relay site overwritten..."))
			} else {
				pterm.Println()
				pterm.Error.Println("Failed to confirm relay site overwrite action")
				os.Exit(1)
			}
		} else {
			pterm.Println()
			pterm.Error.Println("Failed to confirm relay site overwrite action")
			os.Exit(1)
		}
	} else {
		spinner.Info("Relay site not found continuing with installation...")
	}
}

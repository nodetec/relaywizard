package network

import (
	"fmt"
	"github.com/nodetec/relaywiz/pkg/utils"
	"log"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
)

// Function to get SSL certificates using Certbot
func GetCertificates(domainName, email string) bool {

	options := []string{"yes", "no"}

	prompt := pterm.DefaultInteractiveContinue.WithOptions(options)

  pterm.Println()
	pterm.Println(pterm.Cyan("Do you want to obtain SSL certificates using Certbot?"))
	pterm.Println(pterm.Cyan("This step requires that you already have a configured domain name."))
	pterm.Println(pterm.Cyan("You can always re-run this installer after you have configured your domain name."))
  pterm.Println()

	result, _ := prompt.Show()

	if result == "no" {
		return false
	}

	pterm.Println()

	spinner, _ := pterm.DefaultSpinner.Start("Checking SSL certificates...")

	dirName := utils.GetDirectoryName(domainName)

	// Check if certificates already exist
	if utils.FileExists(fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", domainName)) &&
		utils.FileExists(fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", domainName)) {
		spinner.Info("SSL certificates already exist.")
		return true
	}

	err := os.MkdirAll(fmt.Sprintf("/var/www/%s/.well-known/acme-challenge/", dirName), 0755)
	if err != nil {
		log.Fatalf("Error creating directories for Certbot: %v", err)
	}

	spinner.UpdateText("Obtaining SSL certificates...")
	cmd := exec.Command("certbot", "certonly", "--webroot", "-w", fmt.Sprintf("/var/www/%s", dirName), "-d", domainName, "--email", email, "--agree-tos", "--no-eff-email", "-q")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Certbot failed to obtain the certificate for %s: %v", domainName, err)
	}

	spinner.Success("SSL certificates obtained successfully.")
	return true
}

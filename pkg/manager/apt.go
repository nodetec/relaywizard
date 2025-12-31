package manager

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

// Function to check if a package is installed
func isPackageInstalled(currentUsername, packageName string) bool {
	out, err := exec.Command("dpkg-query", "-W", "-f='${Status}'", packageName).Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			errorCode := exitError.ExitCode()
			// Package not found
			if errorCode == 1 {
				return false
			} else {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to check if package is installed: %v", err))
				pterm.Println()
				pterm.Error.Printfln("Failed to check if package is installed: %v", err)
				os.Exit(1)
			}
		}
	}

	status := string(out)

	if status == "'install ok installed'" {
		return true
	}

	return false
}

// Function to install necessary packages
func AptInstallPackages(selectedRelayOption, currentUsername string) {
	spinner, _ := pterm.DefaultSpinner.Start("Updating and installing packages...")

	if currentUsername == relays.RootUser {
		err := exec.Command("apt", "update", "-qq").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to update packages: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to update packages: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "apt", "update", "-qq").Run()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to update packages: %v", err))
			pterm.Println()
			pterm.Error.Printfln("Failed to update packages: %v", err)
			os.Exit(1)
		}
	}

	packages := []string{"sysvinit-utils", "lsof", "curl", "gnupg", "openssh-server", "ufw", "fail2ban", "nginx", "certbot", "python3-certbot-nginx"}

	testedPackageVersions := []string{SysvinitUtilsTestedVersion, LsofTestedVersion, CurlTestedVersion, GnuPGTestedVersion, OpenSSHServerTestedVersion, UfwTestedVersion, Fail2banTestedVersion, NginxTestedVersion, CertbotTestedVersion, Python3CertbotNginxTestedVersion}

	if selectedRelayOption == relays.NostrRsRelayName || selectedRelayOption == relays.StrfryRelayName || selectedRelayOption == relays.WotRelayName || selectedRelayOption == relays.Strfry29RelayName {
		packages = append(packages, "git")
		testedPackageVersions = append(testedPackageVersions, GitTestedVersion)
	}

	if selectedRelayOption == relays.NostrRsRelayName {
		packages = append(packages, "sqlite3", "libsqlite3-dev")
		testedPackageVersions = append(testedPackageVersions, SQLite3TestedVersion, LibSQLite3DevTestedVersion)
	}

	if selectedRelayOption == relays.StrfryRelayName || selectedRelayOption == relays.Strfry29RelayName {
		packages = append(packages, "libssl-dev", "zlib1g-dev", "liblmdb-dev", "libflatbuffers-dev", "libsecp256k1-dev", "libzstd-dev")
		testedPackageVersions = append(testedPackageVersions, LibSSLDevTestedVersion, Zlib1gDevTestedVersion, LibLMDBDevTestedVersion, LibFlatBuffersDevTestedVersion, Libsecp256k1DevTestedVersion, LibzstdDevTestedVersion)
	}

	for index, p := range packages {
		// Check if package is installed, install if not
		if isPackageInstalled(currentUsername, p) {
			spinner.UpdateText(fmt.Sprintf("%s is already installed.", p))
		} else {
			spinner.UpdateText(fmt.Sprintf("Installing %s...", p))
			if currentUsername == relays.RootUser {
				err := exec.Command("apt", "install", "-y", "-qq", p).Run()
				if err != nil {
					logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to install packages: %v", err))
					pterm.Println()
					pterm.Error.Printfln("Failed to install packages: %v", err)
					os.Exit(1)
				}
			} else {
				err := exec.Command("sudo", "apt", "install", "-y", "-qq", p).Run()
				if err != nil {
					logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to install packages: %v", err))
					pterm.Println()
					pterm.Error.Printfln("Failed to install packages: %v", err)
					os.Exit(1)
				}
			}
		}

		// Check if the installed package version matches the tested version, if they don't match then log a warning to the rwz log file
		dpkgSOutput, err := exec.Command("dpkg", "-s", p).CombinedOutput()
		if err != nil {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to check the version for the %s package: %v", p, err))
			pterm.Println()
			pterm.Error.Printfln("Failed to check the version for the %s package: %v", p, err)
			os.Exit(1)
		}

		dpkgStatusOutput := string(dpkgSOutput)

		_, _, found := strings.Cut(dpkgStatusOutput, testedPackageVersions[index])

		if !found {
			logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Warning: The installed version of %s is different than the tested version", p))
		}
	}

	spinner.Success("Packages updated and installed")
}

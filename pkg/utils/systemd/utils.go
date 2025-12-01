package systemd

import (
	"os"
	"os/exec"
	"text/template"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

type ServiceFileParams struct {
	RelayUser      string
	EnvFilePath    string
	BinaryFilePath string
	ConfigFilePath string
}

func CreateServiceFile(currentUsername, serviceFilePath, serviceTemplate, serviceFilePermissionsAsString string, serviceFileParams *ServiceFileParams) {
	if currentUsername == relays.RootUser {
		serviceFile, err := os.Create(serviceFilePath)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create service file: %v", err)
			os.Exit(1)
		}
		defer serviceFile.Close()

		serviceTmpl, err := template.New("service").Parse(serviceTemplate)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to parse service template: %v", err)
			os.Exit(1)
		}

		err = serviceTmpl.Execute(serviceFile, struct{ RelayUser, EnvFilePath, BinaryFilePath, ConfigFilePath string }{RelayUser: serviceFileParams.RelayUser, EnvFilePath: serviceFileParams.EnvFilePath, BinaryFilePath: serviceFileParams.BinaryFilePath, ConfigFilePath: serviceFileParams.ConfigFilePath})
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to execute service template: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "touch", serviceFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create service file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, serviceFilePath, "0666")

		serviceFile, err := os.OpenFile(serviceFilePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open service file: %v", err)
			os.Exit(1)
		}
		defer serviceFile.Close()

		serviceTmpl, err := template.New("service").Parse(serviceTemplate)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to parse service template: %v", err)
			os.Exit(1)
		}

		err = serviceTmpl.Execute(serviceFile, struct{ RelayUser, EnvFilePath, BinaryFilePath, ConfigFilePath string }{RelayUser: serviceFileParams.RelayUser, EnvFilePath: serviceFileParams.EnvFilePath, BinaryFilePath: serviceFileParams.BinaryFilePath, ConfigFilePath: serviceFileParams.ConfigFilePath})
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to execute service template: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, serviceFilePath, serviceFilePermissionsAsString)
	}
}

func Reload(currentUsername string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "daemon-reload").Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to reload systemd daemon: %v", err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "daemon-reload").Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to reload systemd daemon: %v", err)
			os.Exit(1)
		}
	}
}

func EnableService(currentUsername, name string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "enable", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to enable %s service: %v", name, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "enable", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to enable %s service: %v", name, err)
			os.Exit(1)
		}
	}
}

func StartService(currentUsername, name string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "start", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to start %s service: %v", name, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "start", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to start %s service: %v", name, err)
			os.Exit(1)
		}
	}
}

func DisableService(currentUsername, name string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "disable", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to disable %s service: %v", name, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "disable", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to disable %s service: %v", name, err)
			os.Exit(1)
		}
	}
}

func StopService(currentUsername, name string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "stop", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to stop %s service: %v", name, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "stop", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to stop %s service: %v", name, err)
			os.Exit(1)
		}
	}
}

func RestartService(currentUsername, name string) {
	if currentUsername == relays.RootUser {
		err := exec.Command("systemctl", "restart", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to restart %s service: %v", name, err)
			os.Exit(1)
		}
	} else {
		err := exec.Command("sudo", "systemctl", "restart", name).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to restart %s service: %v", name, err)
			os.Exit(1)
		}
	}
}

func DisableAndStopService(currentUsername, path, name string) {
	if files.FileExists(path) {
		DisableService(currentUsername, name)
		StopService(currentUsername, name)
	}
}

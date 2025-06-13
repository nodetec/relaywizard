package systemd

import (
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"text/template"
)

type ServiceFileParams struct {
	RelayUser      string
	EnvFilePath    string
	BinaryFilePath string
	ConfigFilePath string
}

func CreateServiceFile(serviceFilePath, serviceTemplate string, serviceFileParams *ServiceFileParams) {
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
}

func Reload() {
	err := exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to reload systemd daemon: %v", err)
		os.Exit(1)
	}
}

func EnableService(name string) {
	err := exec.Command("systemctl", "enable", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to enable %s service: %v", name, err)
		os.Exit(1)
	}
}

func StartService(name string) {
	err := exec.Command("systemctl", "start", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to start %s service: %v", name, err)
		os.Exit(1)
	}
}

func DisableService(name string) {
	err := exec.Command("systemctl", "disable", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to disable %s service: %v", name, err)
		os.Exit(1)
	}
}

func StopService(name string) {
	err := exec.Command("systemctl", "stop", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to stop %s service: %v", name, err)
		os.Exit(1)
	}
}

func ReloadService(name string) {
	err := exec.Command("systemctl", "reload", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to reload %s service: %v", name, err)
		os.Exit(1)
	}
}

func RestartService(name string) {
	err := exec.Command("systemctl", "restart", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Printfln("Failed to restart %s service: %v", name, err)
		os.Exit(1)
	}
}

func DisableAndStopService(path, name string) {
	if files.FileExists(path) {
		DisableService(name)
		StopService(name)
	}
}
